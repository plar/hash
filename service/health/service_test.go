package health_test

import (
	"testing"
	"time"

	"github.com/plar/hash/service/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HealthServiceSuit struct {
	suite.Suite
	svc health.Service
}

func (s *HealthServiceSuit) SetupTest() {
	s.svc = health.NewService()
}
func (s *HealthServiceSuit) TestDefaults() {
	assert := assert.New(s.T())
	assert.False(s.svc.IsHealthy())
	now := time.Now().Unix()
	assert.GreaterOrEqual(s.svc.Now(), now)
}
func (s *HealthServiceSuit) TestHealthy() {
	assert := assert.New(s.T())
	s.svc.Healthy()
	assert.True(s.svc.IsHealthy())
}
func (s *HealthServiceSuit) TestUnhealthy() {
	assert := assert.New(s.T())
	s.svc.Unhealthy()
	assert.False(s.svc.IsHealthy())
}
func (s *HealthServiceSuit) TestNow() {
	assert := assert.New(s.T())
	now := time.Now().Unix()
	assert.GreaterOrEqual(s.svc.Now(), now)
}
func TestHealthHandlerSuit(t *testing.T) {
	suite.Run(t, new(HealthServiceSuit))
}
