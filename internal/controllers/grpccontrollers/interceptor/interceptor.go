package interceptor

import (
	"context"
	"errors"

	"github.com/EugeneTsydenov/chesshub-user-service/internal/controllers/grpccontrollers/grpcerrors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ErrorHandlingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			requestID := ctx.Value("request-id").(string)
			logger.
				WithField("server", info.Server).
				WithField("method", info.FullMethod).
				WithField("request-id", requestID).
				WithField("cause", errors.Unwrap(err)).
				WithField("error", err).
				Error("[Error handling interceptor]: gRPC request failed")
			return nil, grpcerrors.ToGRPCError(err)
		}

		return resp, nil
	}
}
