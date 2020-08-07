package repository

import (
	"errors"

	"github.com/plar/hash/domain"
)

var ErrHashNotFound = errors.New("hash not found")

// HashRepository represents a persistence layer.
type HashRepository interface {
	// NewID generates a new HashID.
	NewID() domain.HashID

	// Load loads a password hash from the repository.
	// If the repository does not contain a password hash then ErrHashNotFound is returned.
	Load(id domain.HashID) (domain.Hash, error)

	// Save saves a new password hash to the repository.
	Save(id domain.HashID, hash []byte)
}
