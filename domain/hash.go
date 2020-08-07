package domain

import (
	"encoding/base64"
	"fmt"
)

// HashID is the hash identifier
type HashID int64

type Hash struct {
	ID   HashID
	Hash []byte
}

func (h HashID) Bytes() []byte {
	return []byte(fmt.Sprintf("%v", h))
}

func (h Hash) Base64() []byte {
	encHash := make([]byte, base64.StdEncoding.EncodedLen(len(h.Hash)))
	base64.StdEncoding.Encode(encHash, h.Hash)
	return encHash
}

func (h Hash) String() string {
	if h.ID == 0 {
		return "hash{}"
	}
	return fmt.Sprintf("hash{ID: %v, Hash: %s}", h.ID, h.Base64())
}
