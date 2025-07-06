package grpcinterceptors

import (
	"context"

	"github.com/EugeneTsydenov/chesshub-user-service/cmd/user/app/tracker"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequestTracking(tracker *tracker.RequestTracker, logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger.Info("Tracking request")

		if tracker.IsShuttingDown() {
			logger.Warnf("Rejected request to %s: service is shutting down", info.FullMethod)
			return nil, status.Error(codes.Unavailable, "Service is shutting down")
		}

		requestID := uuid.New().String()
		metadata := map[string]interface{}{
			"method": info.FullMethod,
		}

		ctx = context.WithValue(ctx, "request-id", requestID)

		tracker.Begin(requestID, metadata)
		resp, err := handler(ctx, req)
		tracker.End(requestID)

		return resp, err
	}
}
