package shared

import (
	"slices"
	"sync"
	"time"
)

// Filter offre un servizio generico di sampling degli rtt relativi a un certo nodo remoto.
type Filter interface {
	Update(node Node, rtt time.Duration) time.Duration
}

// MPFilter filtro non lineare a percentile p.
type MPFilter struct {
	mu         *sync.RWMutex
	windows    map[Node][]time.Duration
	windowSize int
	p          float64
}

func NewMPFilter(windowSize int, p float64) Filter {
	return &MPFilter{mu: &sync.RWMutex{}, windowSize: windowSize, p: p, windows: make(map[Node][]time.Duration)}
}

func (f *MPFilter) Update(node Node, rtt time.Duration) time.Duration {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.windows[node]
	if !ok {
		f.windows[node] = make([]time.Duration, 0)
	}
	window := f.windows[node]
	if len(window) < f.windowSize {
		window = append(window, rtt)
		f.windows[node] = window
		return rtt
	}
	window = append(window[1:], rtt)
	f.windows[node] = window
	samples := make([]time.Duration, f.windowSize)
	copy(samples, window)
	slices.Sort(samples)
	i := int(float64(len(samples)) * f.p)
	return samples[i]
}

type RawFilter struct {
}

func NewRawFilter() *RawFilter {
	return &RawFilter{}
}

func (r *RawFilter) Update(node Node, rtt time.Duration) time.Duration {
	return rtt
}
