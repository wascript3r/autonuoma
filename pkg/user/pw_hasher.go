package user

type PwHasher interface {
	Hash(password string) (string, error)
	Validate(hash, password string) error
}
