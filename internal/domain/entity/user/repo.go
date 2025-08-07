package user

import (
	"context"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/enums"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error

	Search(ctx context.Context, query string, limit, offset int) ([]*User, error)
	GetActiveUsers(ctx context.Context, limit, offset int) ([]*User, error)
	CheckUsernameAvailable(ctx context.Context, username string) (bool, error)
	CheckEmailAvailable(ctx context.Context, email string) (bool, error)
}

type ProfileRepository interface {
	Create(ctx context.Context, user *Profile) (*Profile, error)
	GetByUserID(ctx context.Context, userID int64) (*Profile, error)
	Update(ctx context.Context, profile *Profile) error
	GetPublicProfiles(ctx context.Context, userIDs []uuid.UUID) ([]*Profile, error)
}

type RatingRepository interface {
	GetUserRatings(ctx context.Context, userID uuid.UUID) ([]*Rating, error)
	GetUserRating(ctx context.Context, userID uuid.UUID, timeControl enums.TimeControl) (*Rating, error)
	UpdateRating(ctx context.Context, rating *Rating) error
	CreateRatingHistory(ctx context.Context, history *RatingHistory) error
	GetRatingHistory(ctx context.Context, userID uuid.UUID, timeControl enums.TimeControl, limit int) ([]*RatingHistory, error)
	GetLeaderboard(ctx context.Context, timeControl enums.TimeControl, limit, offset int) ([]*Rating, error)
	GetUserRank(ctx context.Context, userID uuid.UUID, timeControl enums.TimeControl) (int, error)
}
