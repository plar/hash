package hasher_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/plar/hash/domain"
	"github.com/plar/hash/infra/persistence/memory"
	"github.com/plar/hash/server/config"
	"github.com/plar/hash/service/hasher"
	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	config.DefaultConfig
}

func (c *testConfig) TaskDelay() time.Duration {
	return 0
}

func Config() config.Config {
	return &testConfig{}
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	repo := memory.NewHashRepository()

	svc := hasher.New(repo, Config())
	assert.NotNil(svc)
	svc.Stop()
}

func TestCreate(t *testing.T) {
	var err error
	assert := assert.New(t)

	repo := memory.NewHashRepository()

	svc := hasher.New(repo, Config())

	// bad password
	_, err = svc.Create("")
	assert.Error(err)
	assert.True(errors.Is(err, hasher.ErrInvalidPassword))

	// good password
	hashID, err := svc.Create("angryMonkey")
	assert.NoError(err)
	assert.Equal(domain.HashID(1), hashID)

	var wg sync.WaitGroup
	var hash domain.Hash

	wg.Add(1)
	go func() {
		for {
			hash, err = repo.Load(hashID)
			if err == nil {
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()

	assert.Equal("hash{ID: 1, Hash: ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==}", hash.String())
	assert.NoError(err)

	svc.Stop()
}

func TestGet(t *testing.T) {
	var err error
	assert := assert.New(t)

	repo := memory.NewHashRepository()

	svc := hasher.New(repo, Config())

	// bad password
	_, err = svc.Create("")
	assert.Error(err)
	assert.True(errors.Is(err, hasher.ErrInvalidPassword))

	// good password
	hashID, err := svc.Create("password")
	assert.NoError(err)
	assert.Equal(domain.HashID(1), hashID)

	var wg sync.WaitGroup
	var hash domain.Hash

	wg.Add(1)
	go func() {
		for {
			hash, err = svc.Get(hashID)
			if err == nil {
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()

	assert.Equal("hash{ID: 1, Hash: sQnzu7wkTrgkQZF+0G1hi5AI3Qmzvv0bXgc5THBqi7mAsdd4Xll27ASbRt9fEyavWi6m0QP9B8lThf+rDKy8hg==}", hash.String())
	assert.NoError(err)

	// bad id
	hash, err = svc.Get(domain.HashID(2))
	assert.Equal(domain.Hash{}, hash)
	assert.Error(err)

	svc.Stop()
}
