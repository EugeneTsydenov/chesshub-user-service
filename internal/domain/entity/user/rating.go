package user

import (
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/enums"
	domainerrors "github.com/EugeneTsydenov/chesshub-user-service/internal/domain/errors"
	"github.com/google/uuid"
	"time"
)

const (
	defaultRating = 1200
)

type Rating struct {
	Id           uuid.UUID
	UserID       int64
	TimeControl  enums.TimeControl
	Rating       int
	PeakRating   int
	LowestRating int
	GamesPlayed  int
	Wins         int
	Losses       int
	Draws        int
	LastGameAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (r *Rating) Initialize(userID int64, tc enums.TimeControl) error {
	err := r.GenerateID()
	if err != nil {
		return err
	}

	r.UserID = userID
	r.TimeControl = tc

	r.Rating = defaultRating
	r.PeakRating = defaultRating
	r.LowestRating = defaultRating

	return nil
}

func (r *Rating) GenerateID() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return domainerrors.ErrGeneratingSessionID
	}

	r.Id = id

	return nil
}
