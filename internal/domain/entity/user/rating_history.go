package user

import (
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/enums"
	"github.com/google/uuid"
	"time"
)

type RatingHistory struct {
	Id           uuid.UUID
	UserID       int64
	TimeControl  enums.TimeControl
	OldRating    int
	NewRating    int
	RatingRange  int
	GameID       *uuid.UUID
	ChangeReason enums.ChangeReason
	CreatedAt    time.Time
}
