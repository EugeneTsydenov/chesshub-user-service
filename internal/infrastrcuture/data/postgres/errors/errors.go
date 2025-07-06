package errors

import (
	"context"
	"errors"
	"fmt"
)

type ErrMapper func(err error) error

func WrapWithMapper(op string, err error, mapper ErrMapper) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%s: context error: %w", op, err)
	}

	if mapper != nil {
		return fmt.Errorf("%s: %w", op, mapper(err))
	}

	return fmt.Errorf("%s: unexpected error: %w", op, err)
}
