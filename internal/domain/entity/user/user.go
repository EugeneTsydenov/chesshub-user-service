package user

import (
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/enums"
	"time"

	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/email"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/password"
	"github.com/EugeneTsydenov/chesshub-user-service/internal/domain/valueobjects/tag"
)

type User struct {
	id              int64
	tag             *tag.Tag
	email           *email.Email
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

func (u *User) ID() int64 {
	return u.id
}

func (u *User) Tag() *tag.Tag {
	return u.tag
}

func (u *User) Email() *email.Email {
	return u.email
}

func (u *User) Password() *password.HashedPassword {
	return u.password
}

func (u *User) Status() enums.UserStatus {
	return u.status
}

func (u *User) IsVerified() bool {
	return u.isVerified
}

func (u *User) IsPremium() bool {
	return u.isPremium
}

func (u *User) EmailVerifiedAt() *time.Time {
	return u.emailVerifiedAt
}

func (u *User) PremiumUntil() *time.Time {
	return u.premiumUntil
}

func (u *User) Language() string {
	return u.language
}

func (u *User) LastActiveAt() time.Time {
	return u.lastActiveAt
}

func (u *User) LastLoginAt() *time.Time {
	return u.lastLoginAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) Initialize() {
	u.RefreshLastActiveAt()

	u.status = enums.UserStatusActive

	u.isVerified = false
	u.isPremium = false
}

func (u *User) RefreshLastActiveAt() {
	u.lastActiveAt = time.Now()
}
