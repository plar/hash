package health

import (
	"sync/atomic"
	"time"
)

// Service interface declares health service methods.
type Service interface {
	Healthy()
	Unhealthy()
	IsHealthy() bool
	Now() int64
}
type service struct {
	healthy int32
}

var _ Service = &service{}

func (s *service) Healthy() {
	atomic.StoreInt32(&s.healthy, 1)
}
func (s *service) Unhealthy() {
	atomic.StoreInt32(&s.healthy, 0)
}
func (s *service) IsHealthy() bool {
	return atomic.LoadInt32(&s.healthy) == 1
}
func (s *service) Now() int64 {
	return time.Now().Unix()
}

// NewService creates a health service.
func NewService() Service {
	return &service{}
}
