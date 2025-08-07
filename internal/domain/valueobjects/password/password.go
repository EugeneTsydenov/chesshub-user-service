package password

import (
	"fmt"
	"unicode"
)

const (
	minLen = 8
	maxLen = 64
)

type Hasher interface {
	Hash(value string) (string, error)
	Compare(hash string, plain string) error
}

type PlainPassword struct {
	value string
}

func NewPlainPassword(value string) (*PlainPassword, error) {
	if err := validate(value); err != nil {
		return nil, err
	}

	return &PlainPassword{value}, nil
}

func validate(value string) error {
	if len(value) < minLen || len(value) > maxLen {
		return fmt.Errorf("password length must be between %d and %d", minLen, maxLen)
	}

	var (
		hasLower   bool
		hasUpper   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, c := range value {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}

	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func (p *PlainPassword) Value() string {
	return p.value
}

func (p *PlainPassword) String() string {
	return p.value
}

func (p *PlainPassword) Hash(hasher Hasher) (*HashedPassword, error) {
	hashed, err := hasher.Hash(p.value)
	if err != nil {
		return nil, err
	}

	return &HashedPassword{value: hashed}, nil
}

type HashedPassword struct {
	value string
}

func NewHashedPassword(raw string) *HashedPassword {
	return &HashedPassword{raw}
}

func (p *HashedPassword) Value() string {
	return p.value
}

func (p *HashedPassword) String() string {
	return p.value
}

func (p *HashedPassword) Compare(plain string, hasher Hasher) error {
	return hasher.Compare(p.value, plain)
}
