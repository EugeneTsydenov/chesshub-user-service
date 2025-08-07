package services

import (
	"context"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/entity/user"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/interfaces"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/password"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/tag"
)

type UserService struct {
	hasher password.Hasher

	repo        user.Repository
	profileRepo user.ProfileRepository
}

func NewUserService(hasher password.Hasher, repo user.Repository, profileRepository user.ProfileRepository) *UserService {
	return &UserService{
		hasher:      hasher,
		repo:        repo,
		profileRepo: profileRepository,
	}
}

func (s *UserService) InitializeUserData(ctx context.Context, user *user.User, profile *user.Profile) (*user.User, *user.Profile, error) {
	user.Initialize()

	profile.Initialize(u.ID())

	return u, p, nil
}

//func (s *UserService) SecureUser(u *user.User, plainPassword password.PlainPassword) error {
//	hash, err := s.hasher.Hash(u.Password())
//	if err != nil {
//		return err
//	}
//
//	u.HashPassword(hash)
//
//	return nil
//}

func BuildUserTag(rawTag *string) (*tag.Tag, error) {
	if rawTag == nil {
		generatedTag, err := tag.Generate()
		if err != nil {
			return nil, err
		}

		return generatedTag, nil
	}

	parsedTag, err := tag.New(*rawTag)
	if err != nil {
		return nil, err
	}

	return parsedTag, nil
}
