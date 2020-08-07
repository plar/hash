package memory

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/plar/hash/domain"
	"github.com/plar/hash/domain/repository"
)

// Lets use a simple memory storage for hashes with multiple buckets to avoid locks collions.

const defTotalBuckets = 8

type bucket struct {
	sync.RWMutex
	storage map[domain.HashID][]byte
}

type hashRepository struct {
	totalBuckets int
	buckets      []bucket
	curID        int64
}

var _ repository.HashRepository = &hashRepository{}

// NewHashRepository creates a new hash repository
func NewHashRepository() repository.HashRepository {
	return NewHashRepositoryWithBuckets(0) // use default
}

// NewHashRepositoryWithBuckets creates a new hash repository with the custom number of buckets
func NewHashRepositoryWithBuckets(totalBuckets int) repository.HashRepository {
	if totalBuckets <= 0 {
		totalBuckets = defTotalBuckets
	}

	buckets := make([]bucket, totalBuckets)
	for i := 0; i < totalBuckets; i++ {
		buckets[i] = bucket{storage: make(map[domain.HashID][]byte, 0)}
	}

	return &hashRepository{
		totalBuckets: totalBuckets,
		buckets:      buckets,
	}
}

// NewID generates a new HashID.
func (r *hashRepository) NewID() domain.HashID {
	return domain.HashID(atomic.AddInt64(&r.curID, 1))
}

// Load loads a password hash from the repository.
// If the repository does not contain a password hash then repository.ErrHashNotFound error is returned.
func (r *hashRepository) Load(id domain.HashID) (hash domain.Hash, err error) {
	bid := int(id) % r.totalBuckets

	r.buckets[bid].RLock()
	hash, err = r.loadFromBucket(bid, id)
	r.buckets[bid].RUnlock()
	return
}

func (r *hashRepository) loadFromBucket(bid int, id domain.HashID) (domain.Hash, error) {
	hash, ok := r.buckets[bid].storage[id]
	if !ok {
		return domain.Hash{}, fmt.Errorf("HashID '%v': %w", id, repository.ErrHashNotFound)
	}

	return domain.Hash{
		ID:   id,
		Hash: hash,
	}, nil
}

// Save saves a new password hash to the repository.
func (r *hashRepository) Save(id domain.HashID, hash []byte) {
	bid := int(id) % r.totalBuckets

	r.buckets[bid].Lock()
	r.buckets[bid].storage[id] = hash
	r.buckets[bid].Unlock()
}
