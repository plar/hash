package hasher

import (
	"time"

	"github.com/plar/hash/domain"
	"github.com/plar/hash/domain/repository"
	"github.com/plar/hash/infra/pool"
	"github.com/plar/hash/server/config"
)

// Service interface declares the Hasher service methods
type Service interface {
	// Create creates a new hash request
	Create(password string) (domain.HashID, error)

	// Get retrieves a password hash by id
	// returns an error if the password hash was not found
	Get(id domain.HashID) (domain.Hash, error)

	// Stop the server
	Stop()
}

var _ Service = &service{}

type service struct {
	taskQueue  chan pool.Task
	dispatcher *pool.Dispatcher

	hashRepo repository.HashRepository
	delay    time.Duration
}

// New creates a new hasher service
func New(hashRepo repository.HashRepository, cfg config.Config) Service {

	// Create the task queue
	taskQueue := make(chan pool.Task, cfg.QueueSize())

	// Start the dispatcher with 8 workers
	dispatcher := pool.NewDispatcher(taskQueue, cfg.TotalWorkers())
	dispatcher.Run()

	return &service{
		taskQueue:  taskQueue,
		dispatcher: dispatcher,
		hashRepo:   hashRepo,
		delay:      cfg.TaskDelay(),
	}
}

func (s *service) Create(password string) (domain.HashID, error) {
	// pre-validate input args
	if len(password) == 0 {
		return 0, ErrInvalidPassword
	}

	hashID := s.hashRepo.NewID()

	// Execute calculation of the hash code asynchronously.
	// We can get stuck here if we have more then 1000 tasks in the taskQueue.
	// A goroutine can be used before `go s.dispatcher.Dispatch(...)`
	// but in that case we can run out of memory.
	s.dispatcher.Dispatch(pool.Task{
		ID: int64(hashID),
		Handler: func() {
			// sim long-running process...
			time.Sleep(s.delay)

			// calc hash
			encryptor := NewSHA512Encryptor()
			hash := encryptor.Hash([]byte(password))

			// update storage
			s.hashRepo.Save(hashID, hash)
		},
	})

	return hashID, nil
}

func (s *service) Get(id domain.HashID) (domain.Hash, error) {
	hash, err := s.hashRepo.Load(id)
	if err != nil {
		return domain.Hash{}, err
	}

	return hash, nil
}

func (s *service) Stop() {
	s.dispatcher.Stop()
}
