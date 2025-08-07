package user

import (
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/enums"
	"time"

	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/email"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/password"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/publicname"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/tag"
)

type Builder struct {
	id              int64
	tag             *tag.Tag
	email           *email.Email
	publicName      *publicname.PublicName
	password        *password.HashedPassword
	status          enums.UserStatus
	isVerified      bool
	isPremium       bool
	emailVerifiedAt *time.Time
	premiumUntil    *time.Time
	language        string
	lastActiveAt    time.Time
	lastLoginAt     *time.Time
	updatedAt       time.Time
	createdAt       time.Time
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithID(id int64) *Builder {
	b.id = id

	return b
}

func (b *Builder) WithTag(tag *tag.Tag) *Builder {
	b.tag = tag

	return b
}

func (b *Builder) WithEmail(email *email.Email) *Builder {
	b.email = email

	return b
}

func (b *Builder) WithPassword(password *password.HashedPassword) *Builder {
	b.password = password

	return b
}

func (b *Builder) WithStatus(status enums.UserStatus) *Builder {
	b.status = status

	return b
}

func (b *Builder) WithIsVerified(isVerified bool) *Builder {
	b.isVerified = isVerified

	return b
}

func (b *Builder) WithIsPremium(isPremium bool) *Builder {
	b.isPremium = isPremium

	return b
}

func (b *Builder) WithEmailVerifiedAt(emailVerifiedAt *time.Time) *Builder {
	b.emailVerifiedAt = emailVerifiedAt

	return b
}

func (b *Builder) WithPremiumUntil(premiumUntil *time.Time) *Builder {
	b.premiumUntil = premiumUntil

	return b
}

func (b *Builder) WithLanguage(language string) *Builder {
	b.language = language

	return b
}

func (b *Builder) WithLastLoginAt(lastLoginAt *time.Time) *Builder {
	b.lastLoginAt = lastLoginAt

	return b
}

func (b *Builder) WithLastActiveAt(lastActiveAt time.Time) *Builder {
	b.lastActiveAt = lastActiveAt

	return b
}

func (b *Builder) WithUpdatedAt(updatedAt time.Time) *Builder {
	b.updatedAt = updatedAt

	return b
}

func (b *Builder) WithCreatedAt(createdAt time.Time) *Builder {
	b.createdAt = createdAt

	return b
}

func (b *Builder) Build() *User {
	return &User{
		id:              b.id,
		tag:             b.tag,
		email:           b.email,
		password:        b.password,
		status:          b.status,
		isVerified:      b.isVerified,
		isPremium:       b.isPremium,
		emailVerifiedAt: b.emailVerifiedAt,
		premiumUntil:    b.premiumUntil,
		language:        b.language,
		lastActiveAt:    b.lastActiveAt,
		lastLoginAt:     b.lastLoginAt,
		updatedAt:       b.updatedAt,
		createdAt:       b.createdAt,
	}
}
