package metrics

import (
	"sync"
	"time"
)

const (
	CurrentConnections       = "current_connections"
	MaxConcurrentConnections = "max_concurrent_connections"
	TotalConnections         = "total_connections"
)

var (
	metricTypes = map[string]string{
		CurrentConnections:       "gauge",
		MaxConcurrentConnections: "gauge",
		TotalConnections:         "counter",
	}
)

type Metrics struct {
	sync.RWMutex
	metrics     map[string]int
	fileWriters []*fileWriter
}

func New() *Metrics {
	return &Metrics{
		metrics:     make(map[string]int),
		fileWriters: make([]*fileWriter, 0),
	}
}

func (m *Metrics) WithFileWriter(filename string, format string, frequency time.Duration) *Metrics {
	m.Lock()
	defer m.Unlock()
	m.fileWriters = append(m.fileWriters, newFileWriter(m, filename, format, frequency, true))
	return m
}

func (m *Metrics) Set(key string, value int) {
	m.Lock()
	defer m.Unlock()
	m.metrics[key] = value
}

func (m *Metrics) SetIfGreater(key string, value int) {
	m.Lock()
	defer m.Unlock()
	if value > m.metrics[key] {
		m.metrics[key] = value
	}
}

func (m *Metrics) Inc(key string) {
	m.Lock()
	defer m.Unlock()
	m.metrics[key]++
}

func (m *Metrics) Get(key string) int {
	m.RLock()
	defer m.RUnlock()
	return m.metrics[key]
}

func (m *Metrics) GetAll() map[string]int {
	m.RLock()
	defer m.RUnlock()
	// return a copy to avoid concurrent map access, this should be more performant than using a sync.Map
	metricsCopy := make(map[string]int)
	for k, v := range m.metrics {
		metricsCopy[k] = v
	}
	return metricsCopy
}
