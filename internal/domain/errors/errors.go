package errors

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotVerified     = errors.New("user not verified")
	ErrUserSuspended       = errors.New("user suspended")
	ErrUserBanned          = errors.New("user banned")
	ErrRatingNotFound      = errors.New("rating not found")
	ErrInvalidRatingChange = errors.New("invalid rating change")
	ErrUsernameUnavailable = errors.New("username unavailable")
	ErrEmailUnavailable    = errors.New("email unavailable")

	ErrGeneratingSessionID = errors.New("generating session id failed")
)
