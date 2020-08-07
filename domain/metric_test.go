package domain_test

import (
	"testing"
	"time"

	"github.com/plar/hash/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MetricSuit struct {
	suite.Suite
}

func TestMetric(t *testing.T) {
	assert := assert.New(t)

	m := domain.Metric{}
	assert.Equal(int64(0), m.Count())
	assert.Equal(time.Duration(0), m.Average())

	m.AddCount(1)
	m.AddExecTime(time.Duration(100) * time.Second)
	assert.Equal(int64(1), m.Count())
	assert.Equal(time.Duration(100)*time.Second, m.Average())

	m.AddCount(1)
	m.AddExecTime(time.Duration(100) * time.Second)
	assert.Equal(int64(2), m.Count())
	assert.Equal(time.Duration(100)*time.Second, m.Average())

	m.AddCount(1)
	m.AddExecTime(time.Duration(200) * time.Second)
	assert.Equal(int64(3), m.Count())
	assert.Equal(time.Duration(100+100+200)*time.Second/time.Duration(3), m.Average())
}
