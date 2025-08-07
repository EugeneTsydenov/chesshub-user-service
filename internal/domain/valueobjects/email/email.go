package email

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

func New(value string) (*Email, error) {
	if err := validate(value); err != nil {
		return nil, err
	}

	return &Email{value}, nil
}

func validate(value string) error {
	if value == "" {
		return errors.New("email cannot be empty")
	}

	emailLen := 254

	if len(value) > emailLen {
		return errors.New("email too long")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, value)
	if err != nil {
		return errors.New("invalid email format")
	}

	if !matched {
		return errors.New("invalid email format")
	}

	atCount := strings.Count(value, "@")
	if atCount != 1 {
		return errors.New("email must contain exactly one @ symbol")
	}

	parts := strings.Split(value, "@")
	localPart := parts[0]
	domainPart := parts[1]

	if len(localPart) == 0 || len(localPart) > 64 {
		return errors.New("local part must be between 1 and 64 characters")
	}

	if len(domainPart) == 0 || len(domainPart) > 253 {
		return errors.New("domain part must be between 1 and 253 characters")
	}

	if strings.HasPrefix(value, ".") || strings.HasSuffix(value, ".") {
		return errors.New("email cannot start or end with a dot")
	}

	if strings.Contains(value, "..") {
		return errors.New("email cannot contain consecutive dots")
	}

	return nil
}

func (vo *Email) Value() string {
	return vo.value
}

func (vo *Email) String() string {
	return vo.value
}
