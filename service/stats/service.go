package stats

import (
	"sync"
	"time"

	"github.com/plar/hash/domain"
)

// Service interface declares the Stats service methods
type Service interface {
	TrackMetric(name string, count int64, execTime time.Duration)
	Metric(name string) domain.Metric
}

var _ Service = &service{}

type service struct {
	lock    sync.RWMutex
	metrics map[string]domain.Metric
}

func New() Service {
	return &service{
		metrics: make(map[string]domain.Metric, 0),
	}
}

func (s *service) TrackMetric(name string, count int64, execTime time.Duration) {
	s.lock.Lock()

	m, _ := s.metrics[name]
	m.AddCount(count)
	m.AddExecTime(execTime)
	s.metrics[name] = m

	s.lock.Unlock()
}

func (s *service) Metric(name string) domain.Metric {
	s.lock.RLock()

	m, _ := s.metrics[name]

	s.lock.RUnlock()
	return m
}
