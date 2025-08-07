package usecase

import "context"

type UseCase[Input comparable, Output comparable] interface {
	Execute(ctx context.Context, input Input) (Output, error)
}
