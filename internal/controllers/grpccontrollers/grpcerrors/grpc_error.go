package grpcerrors

import (
	"errors"
	"fmt"

	apperrors "github.com/EugeneTsydenov/chesshub-user-service/internal/app/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPCError(err error) error {
	var appErr *apperrors.AppError

	if errors.As(err, &appErr) {
		return appErrorToGRPCError(appErr)
	}

	return status.Error(codes.Unknown, "Unexpected server error.")
}

func appErrorToGRPCError(err *apperrors.AppError) error {
	switch err.Type {
	case apperrors.InvalidArgument:
		return withDetails(codes.InvalidArgument, err)
	case apperrors.NotFound:
		return status.Error(codes.NotFound, err.Message)
	case apperrors.Conflict:
		return withDetails(codes.AlreadyExists, err)
	case apperrors.Internal:
		return status.Error(codes.Internal, err.Message)
	case apperrors.Unauthenticated:
		return status.Error(codes.Unauthenticated, err.Message)
	case apperrors.Forbidden:
		return status.Error(codes.PermissionDenied, err.Message)
	case apperrors.Canceled:
		return status.Error(codes.Canceled, err.Message)
	case apperrors.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, err.Message)
	default:
		return status.Error(codes.Unknown, err.Message)
	}
}

func withDetails(code codes.Code, err *apperrors.AppError) error {
	errInfo := &errdetails.ErrorInfo{
		Reason:   err.Type.String(),
		Domain:   "session",
		Metadata: err.Metadata,
	}

	st := status.New(code, err.Message)
	detailedStatus, detailErr := st.WithDetails(errInfo)
	if detailErr != nil {
		return fmt.Errorf("st.WithDetails: %w", detailErr)
	}
	return detailedStatus.Err()
}
