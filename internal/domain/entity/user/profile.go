package user

import (
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/publicname"
	"time"
)

type Profile struct {
	UserID            int64
	PublicName        *publicname.PublicName
	Bio               *string
	CountryCode       *string
	City              *string
	BirthDate         *time.Time
	AvatarURL         *string
	CoverImageURL     *string
	IsPublic          bool
	ShowCountry       bool
	WebsiteURL        *string
	TwitchUsername    *string
	YoutubeChannelURL *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (p *Profile) Initialize(userID int64) {
	p.UserID = userID

	p.IsPublic = true
	p.ShowCountry = true
}
