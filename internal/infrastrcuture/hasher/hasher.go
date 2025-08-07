package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

func New(cost int) *Hasher {
	return &Hasher{cost}
}

func (h *Hasher) Hash(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), h.cost)
	return string(bytes), err
}

func (h *Hasher) Compare(hash string, plain string, cost int) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return err
	}

	return nil
}
