package hash

import "github.com/alexedwards/argon2id"

type Hasher interface {
	CreateHash(string) (string, error)
	CompareHash(string, string) (bool, error)
}

type Hash struct {
	params *argon2id.Params
}

func NewHash() *Hash {
	return &Hash{
		params: argon2id.DefaultParams,
	}
}

func (hash *Hash) CreateHash(password string) (string, error) {
	h, err := argon2id.CreateHash(password, hash.params)
	if err != nil {
		return "", err
	}

	return h, nil
}

func (hash *Hash) CompareHash(password, hashed string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hashed)
	if err != nil {
		return false, err
	}

	return match, nil
}
