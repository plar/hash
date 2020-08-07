package memory

import (
	"errors"
	"testing"

	"github.com/plar/hash/domain"
	"github.com/plar/hash/domain/repository"
	"github.com/stretchr/testify/assert"
)

func TestNewHashRepository(t *testing.T) {
	assert := assert.New(t)

	r := NewHashRepository()
	ri, ok := r.(*hashRepository)
	assert.True(ok)
	assert.Equal(defTotalBuckets, ri.totalBuckets)
	assert.Len(ri.buckets, defTotalBuckets)
	assert.Equal(int64(0), ri.curID)

	for _, b := range ri.buckets {
		assert.NotNil(b)
		assert.NotNil(b.storage)
	}
}

func TestNewID(t *testing.T) {
	assert := assert.New(t)

	r := NewHashRepository()
	ri, _ := r.(*hashRepository)

	assert.Equal(domain.HashID(1), r.NewID())
	assert.Equal(domain.HashID(2), r.NewID())
	assert.Equal(domain.HashID(3), r.NewID())
	assert.Equal(int64(3), ri.curID)
}

func TestLoadAndSave(t *testing.T) {
	assert := assert.New(t)

	r := NewHashRepositoryWithBuckets(8)
	ri, _ := r.(*hashRepository)
	assert.NotNil(ri)

	for i := 0; i < 8+4; i++ {
		id := domain.HashID(i)
		hash := []byte{byte(i)}
		r.Save(id, hash)
		h, err := r.Load(id)
		assert.NoError(err)
		assert.Equal(domain.Hash{id, hash}, h)
	}

	// test bucket distribution
	for i := 0; i < 8+4; i++ {
		bid := i % 8
		id := domain.HashID(i)
		hash := []byte{byte(i)}
		assert.Contains(ri.buckets[bid].storage, id)
		shash := ri.buckets[bid].storage[id]
		assert.Equal(hash, shash)
	}

	id := domain.HashID(0xdead)
	h, err := r.Load(id)
	assert.Equal(domain.Hash{}, h)
	assert.Error(err)
	assert.True(errors.Is(err, repository.ErrHashNotFound))
}

func BenchmarkHashRepositorySave(b *testing.B) {
	r := NewHashRepositoryWithBuckets(8)
	for n := 0; n < b.N; n++ {
		id := domain.HashID(n)
		hash := []byte{byte(n)}
		r.Save(id, hash)
	}
}

func BenchmarkHashRepositoryLoad(b *testing.B) {
	r := NewHashRepositoryWithBuckets(8)
	for n := 0; n < b.N; n++ {
		id := domain.HashID(n)
		hash := []byte{byte(n)}
		r.Save(id, hash)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		id := domain.HashID(n)
		r.Load(id)
	}
}
