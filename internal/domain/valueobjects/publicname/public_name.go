package publicname

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	minLen = 3
	maxLen = 15
)

var (
	ErrInvalidPublicNameLen     = fmt.Errorf("public name must be between %d and %d", minLen, maxLen)
	ErrInvalidPublicNamePattern = errors.New("public name must contain only Latin letters, digits, and special characters")
)

type PublicName struct {
	value string
}

func New(value string) (*PublicName, error) {
	if err := validate(value); err != nil {
		return nil, err
	}

	return &PublicName{value}, nil
}

func validate(value string) error {
	if len(value) < minLen || len(value) > maxLen {
		return ErrInvalidPublicNameLen
	}

	r := regexp.MustCompile("^[a-zA-Z0-9._-]+$")

	if !r.MatchString(value) {
		return ErrInvalidPublicNamePattern
	}

	return nil
}

func (vo *PublicName) Value() string {
	return vo.value
}

func (vo *PublicName) String() string {
	return vo.value
}
