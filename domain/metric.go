package domain

import "time"

type Metric struct {
	count         int64
	totalExecTime time.Duration
}

func (m *Metric) AddCount(c int64) {
	m.count += c
}

func (m *Metric) AddExecTime(t time.Duration) {
	m.totalExecTime += t
}

func (m Metric) Count() int64 {
	return m.count
}

func (m Metric) Average() time.Duration {
	if m.count == 0 {
		return 0
	}

	return (m.totalExecTime / time.Duration(m.count))
}
