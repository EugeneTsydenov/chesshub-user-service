package tag

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/EugeneTsydenov/chesshub-user-service/internal/pkg/taggen"
)

const (
	tagGenLen = 6

	genTagRetries = 3

	minLen = 3
	maxLen = 10
)

type Tag struct {
	value string
}

func New(value string) (*Tag, error) {
	if err := validate(value); err != nil {
		return nil, err
	}

	return &Tag{value}, nil
}

func validate(value string) error {
	if len(value) < minLen || len(value) > maxLen {
		return fmt.Errorf("tag value must be between %d and %d", minLen, maxLen)
	}

	r := regexp.MustCompile("^[a-zA-Z0-9._-]+$")

	if !r.MatchString(value) {
		return errors.New("tag value must contain only alphanumeric characters")
	}

	return nil
}

func Generate() (*Tag, error) {
	for i := 0; i < genTagRetries; i++ {
		tag, err := taggen.Generate(tagGenLen)
		if err == nil {
			return &Tag{tag}, nil
		}

		log.Printf("failed to generate tag, attempt %d: %v", i+1, err)
	}

	return nil, fmt.Errorf("failed to generate tag after %d attempts", genTagRetries)
}

func (vo *Tag) Value() string {
	return vo.value
}

func (vo *Tag) String() string {
	return vo.value
}
