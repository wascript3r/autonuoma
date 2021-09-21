package pwhasher

import "github.com/wascript3r/gocipher/bcrypt"

type PwHasher struct {
	cost int
}

func New(cost int) PwHasher {
	return PwHasher{cost}
}

func (p PwHasher) Hash(password string) (string, error) {
	bs, err := bcrypt.Compute([]byte(password), p.cost)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (p PwHasher) Validate(hash, password string) error {
	return bcrypt.Validate([]byte(hash), []byte(password))
}
