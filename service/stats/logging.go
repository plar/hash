package stats

import (
	"log"
	"time"

	"github.com/plar/hash/domain"
)

type loggingService struct {
	next Service
}

var _ Service = &loggingService{}

func NewLoggingService(next Service) Service {
	return &loggingService{
		next: next,
	}
}

func (s *loggingService) TrackMetric(name string, count int64, execTime time.Duration) {
	defer func() {
		log.Printf("the stats service method=TrackMetric name=%v, count=%v, execTime=%v", name, count, execTime)
	}()
	s.next.TrackMetric(name, count, execTime)
}

func (s *loggingService) Metric(name string) (m domain.Metric) {
	defer func() {
		log.Printf("the stats service method=Metric name=%v => %v", name, m)
	}()
	return s.next.Metric(name)
}
