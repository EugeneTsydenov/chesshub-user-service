package usecase

import (
	"context"

	"github.com/EugeneTsydenov/chesshub-user-service/internal/app/dto"
	apperrors "github.com/EugeneTsydenov/chesshub-user-service/internal/app/errors"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/entity/user"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/interfaces"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/email"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/password"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/publicname"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/tag"
)

type (
	RegisterUser UseCase[*dto.RegisterUserInputDTO, *dto.RegisterUserOutputDTO]

	registerUser struct {
		userRepository user.Repository
		hasher         password.Hasher
		userService    user.Service
	}
)

func NewRegisterUser(repository user.Repository, hasher password.Hasher, service user.Service) RegisterUser {
	return &registerUser{
		userRepository: repository,
		hasher:         hasher,
		userService:    service,
	}
}

func (uc *registerUser) Execute(ctx context.Context, input *dto.RegisterUserInputDTO) (*dto.RegisterUserOutputDTO, error) {
	errs := make(map[string]string)

	emailVO, err := email.New(input.Email)
	if err != nil {
		errs["email"] = err.Error()
	}

	plainPasswordVO, err := password.NewPlainPassword(input.Password)
	if err != nil {
		errs["password"] = err.Error()
	}

	if input.Language == "" {
		errs["language"] = "language required"
	}

	publicNameVO, err := publicname.New(input.PublicName)
	if err != nil {
		errs["publicname"] = err.Error()
	}

	if len(errs) > 0 {
		return nil, apperrors.NewInvalidArgumentError("validation failed", errs)
	}

	hashed, err := plainPasswordVO.Hash(uc.hasher)
	if err != nil {
		errs["password"] = err.Error()
	}

	builder := user.NewBuilder()

	u := builder.
		WithEmail(emailVO).
		WithPassword(hashed).
		WithLanguage(input.Language).
		Build()

	profile := &user.Profile{
		PublicName: publicNameVO,
	}

	u.Initialize()

	_, err = uc.userRepo.Create(ctx, u)
	if err != nil {
		return nil, apperrors.FromDomainError(err)
	}

	return &dto.RegisterUserOutputDTO{
		Message: "User successfully registered.",
	}, nil
}

func (uc *registerUser) buildUser(input *dto.RegisterUserInputDTO) (*user.User, error) {

}
