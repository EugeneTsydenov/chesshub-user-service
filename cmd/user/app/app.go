package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EugeneTsydenov/chesshub-user-service/cmd/user/app/grpcinterceptors"
	"github.com/EugeneTsydenov/chesshub-user-service/cmd/user/app/tracker"
	"github.com/EugeneTsydenov/chesshub-user-service/config"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/controllers/grpccontrollers/interceptor"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/infrastrcuture/data/postgres"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/infrastrcuture/data/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var shutdownTimeout = 30 * time.Second

type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

type App struct {
	requestTracker *tracker.RequestTracker

	config *config.Config
	logger *logrus.Logger

	redisDatabase *redis.Database
	database      *postgres.Database

	gRPCServer *grpc.Server

	shutdownCh  chan struct{}
	shutdowners []Shutdowner
}

func New(cfg *config.Config) *App {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetReportCaller(true)

	return &App{
		requestTracker: tracker.NewRequestTracker(logger),
		config:         cfg,
		logger:         logger,
		shutdownCh:     make(chan struct{}),
	}
}

func (a *App) InitDeps(ctx context.Context) error {
	if err := a.initRedisCache(ctx); err != nil {
		a.logger.Error(err)

		return err
	}

	if err := a.initPgDatabase(ctx); err != nil {
		a.logger.Error(err)

		return err
	}

	return nil
}

func (a *App) initRedisCache(ctx context.Context) error {
	c, err := redis.New(ctx, a.config.Redis.ConnStr())
	if err != nil {
		return err
	}

	cmd := c.Client().Ping(ctx)
	if err = cmd.Err(); err != nil {
		return err
	}

	a.redisDatabase = c
	a.RegisterShutdowner(c)

	return nil
}

func (a *App) initPgDatabase(ctx context.Context) error {
	d, err := postgres.New(ctx, a.config.Database.DSN())
	if err != nil {
		return err
	}

	err = d.Pool().Ping(ctx)
	if err != nil {
		return err
	}

	a.database = d
	a.RegisterShutdowner(d)

	return nil
}

func (a *App) SetupGRPCServer() {
	a.gRPCServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			grpcinterceptors.RequestTracking(a.requestTracker, a.logger),
			interceptor.ErrorHandlingInterceptor(a.logger),
		),
	)
	reflection.Register(a.gRPCServer)
}

func (a *App) Run(ctx context.Context) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	p := fmt.Sprintf(":%v", a.config.App.Port)
	listener, err := net.Listen("tcp", p)
	if err != nil { //nolint:wsl
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		a.logger.Info("Starting gRPC server", "port", p)
		if err := a.gRPCServer.Serve(listener); err != nil { //nolint:wsl
			a.logger.Error("gRPC server error", "error", err)
		}
	}()

	a.logger.Info("Application started successfully")

	select {
	case sig := <-sigCh:
		a.logger.Info("Received shutdown signal", "signal", sig)
	case <-a.shutdownCh:
		a.logger.Info("Shutdown requested programmatically")
	}

	return a.Shutdown(ctx)
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Starting graceful shutdown")

	a.requestTracker.SetShuttingDown(true)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	grpcShutdownDone := make(chan struct{})
	go func() {
		a.logger.Info("Shutting down gRPC server")
		a.gRPCServer.GracefulStop()
		close(grpcShutdownDone)
	}()

	select {
	case <-grpcShutdownDone:
		a.logger.Info("gRPC server shut down")
	case <-ctx.Done():
		a.logger.Warn("gRPC server shutdown timed out, forcing stop")
		a.gRPCServer.Stop()
	}

	if err := a.requestTracker.WaitForCompletion(ctx); err != nil {
		a.logger.Error("Timed out waiting for requests to complete", "error", err)
	}

	a.logger.Info("Waiting for active requests to complete")
	for _, shutdowner := range a.shutdowners {
		if err := shutdowner.Shutdown(ctx); err != nil {
			a.logger.Error("Error shutting down component", "error", err)
		}
	}

	a.logger.Info("Graceful shutdown completed")

	return nil
}

func (a *App) RegisterShutdowner(s Shutdowner) {
	a.shutdowners = append(a.shutdowners, s)
}
