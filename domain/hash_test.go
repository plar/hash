package domain_test

import (
	"testing"

	"github.com/plar/hash/domain"
	"github.com/stretchr/testify/assert"
)

func TestHashID(t *testing.T) {
	assert := assert.New(t)
	var hashID domain.HashID
	assert.Equal([]byte("0"), hashID.Bytes())

	hashID = 0xdeadbeaf
	assert.Equal([]byte("3735928495"), hashID.Bytes())
}

func TestHash(t *testing.T) {
	assert := assert.New(t)
	hash := domain.Hash{
		ID:   123,
		Hash: []byte{0xde, 0xad, 0xbe, 0xaf}}

	assert.Equal([]byte("3q2+rw=="), hash.Base64())
	assert.Equal("hash{ID: 123, Hash: 3q2+rw==}", hash.String())

	assert.Equal("hash{}", domain.Hash{}.String())
}
