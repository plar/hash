package hasher

import (
	"crypto/sha512"
	"errors"
	"hash"
)

// ErrInvalidPassword is retured when password is invalid
var ErrInvalidPassword = errors.New("Password is empty or nil")

type Encryptor interface {
	Hash(password []byte) []byte
}

type sha512Encryptor struct {
	hashFn hash.Hash
}

func NewSHA512Encryptor() Encryptor {
	return &sha512Encryptor{
		hashFn: sha512.New(),
	}
}

// Hash generates SHA512 of password
func (e *sha512Encryptor) Hash(password []byte) []byte {
	e.hashFn.Reset()
	e.hashFn.Write(password)
	return e.hashFn.Sum(nil)
}
