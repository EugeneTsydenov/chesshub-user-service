package taggen

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Generate(tagLen int8) (string, error) {
	tag := make([]byte, tagLen)

	for i := range tag {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}

		tag[i] = charset[n.Int64()]
	}

	return string(tag), nil
}
