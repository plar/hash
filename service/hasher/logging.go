package hasher

import (
	"log"

	"github.com/plar/hash/domain"
)

type loggingService struct {
	next Service
}

// NewLoggingService creates a new loggin service
func NewLoggingService(next Service) Service {
	return &loggingService{
		next: next,
	}
}

func (s *loggingService) Create(password string) (id domain.HashID, err error) {
	defer func() {
		log.Printf("the hasher service method=Create => id=%v, err=%v", id, err)
	}()
	return s.next.Create(password)
}

func (s *loggingService) Get(id domain.HashID) (hash domain.Hash, err error) {
	defer func() {
		log.Printf("the hasher service method=Get id=%v => hash=%v, err=%v", id, hash, err)
	}()
	return s.next.Get(id)
}

func (s *loggingService) Stop() {
	log.Println("the hasher service is stopping")
	s.next.Stop()
	log.Println("the hasher service has been stopped")
}
