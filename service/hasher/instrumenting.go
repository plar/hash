package hasher

import (
	"time"

	"github.com/plar/hash/domain"
	"github.com/plar/hash/service/stats"
)

type instrumentingService struct {
	next  Service
	stats stats.Service
}

func NewInstrumentingService(next Service, stats stats.Service) Service {
	return &instrumentingService{
		next:  next,
		stats: stats,
	}
}

func (s *instrumentingService) Create(password string) (domain.HashID, error) {
	defer func(begin time.Time) {
		s.stats.TrackMetric("Hasher.Create", 1, time.Since(begin))
	}(time.Now())
	return s.next.Create(password)
}

func (s *instrumentingService) Get(id domain.HashID) (domain.Hash, error) {
	defer func(begin time.Time) {
		s.stats.TrackMetric("Hasher.Get", 1, time.Since(begin))
	}(time.Now())
	return s.next.Get(id)
}

func (s *instrumentingService) Stop() {
	// nothing to collect here
	s.next.Stop()
}
