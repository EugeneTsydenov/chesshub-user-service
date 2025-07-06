package errors

import (
	"context"
	"errors"
	"fmt"
)

type ErrorType int

const (
	InvalidArgument ErrorType = iota
	NotFound
	Conflict
	Internal
	Unauthenticated
	Forbidden
	Canceled
	DeadlineExceeded
)

func (c ErrorType) String() string {
	switch c {
	case InvalidArgument:
		return "INVALID_ARGUMENT_ERROR"
	case NotFound:
		return "NOT_FOUND_ERROR"
	case Conflict:
		return "CONFLICT_ERROR"
	case Internal:
		return "INTERNAL_ERROR"
	case Unauthenticated:
		return "UNAUTHENTICATED_ERROR"
	case Forbidden:
		return "FORBIDDEN_ERROR"
	case DeadlineExceeded:
		return "DEADLINE_EXCEEDED_ERROR"
	case Canceled:
		return "CANCELED_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

type metadata map[string]string

type AppError struct {
	Type     ErrorType
	Message  string
	Metadata metadata
	Cause    error
}

func NewAppError(t ErrorType, msg string, metadata metadata, cause error) *AppError {
	return &AppError{
		Type:     t,
		Message:  msg,
		Metadata: metadata,
		Cause:    cause,
	}
}

func (e *AppError) WithMetadata(key string, value string) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(metadata)
	}

	e.Metadata[key] = value
	return e
}

func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%v: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) Join() error {
	return errors.Join(e, e.Cause)
}

func NewInternalError(msg string) *AppError {
	return NewAppError(Internal, msg, nil, nil)
}

func NewInvalidArgumentError(msg string, metadata metadata) *AppError {
	return NewAppError(InvalidArgument, msg, metadata, nil)
}

func NewNotFoundError(msg string) *AppError {
	return NewAppError(NotFound, msg, nil, nil)
}

func NewConflictError(msg string) *AppError {
	return NewAppError(Conflict, msg, nil, nil)
}

func NewUnauthenticatedError(msg string) *AppError {
	return NewAppError(Unauthenticated, msg, nil, nil)
}

func NewForbiddenError(msg string) *AppError {
	return NewAppError(Forbidden, msg, nil, nil)
}

func NewDeadlineExceededError(msg string) *AppError {
	return NewAppError(DeadlineExceeded, msg, nil, nil)
}

func NewCanceledError(msg string) *AppError {
	return NewAppError(Canceled, msg, nil, nil)
}

func FromDomainError(err error) *AppError {
	switch {
	case errors.Is(err, context.Canceled):
		return NewCanceledError("Operation was canceled.").WithCause(err)
	case errors.Is(err, context.DeadlineExceeded):
		return NewDeadlineExceededError("Operation time out.").WithCause(err)
	default:
		return NewInternalError("Unexpected server error.").WithCause(err)
	}
}
