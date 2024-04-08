package vivaldi

import (
	"b4/shared"
	"slices"
	"time"
)

// Filter offre un servizio generico di sampling degli rtt relativi ad un certo nodo remoto.
type Filter interface {
	Update(node shared.Node, rtt time.Duration) time.Duration
}

// MPFilter filtro non lineare a percentile p.
type MPFilter struct {
	windows    map[shared.Node][]time.Duration
	windowSize int
	p          float64
}

func NewMPFilter(windowSize int, p float64) Filter {
	return &MPFilter{windowSize: windowSize, p: p, windows: make(map[shared.Node][]time.Duration)}
}

func (M *MPFilter) Update(node shared.Node, rtt time.Duration) time.Duration {
	_, ok := M.windows[node]
	if !ok {
		M.windows[node] = make([]time.Duration, 0)
	}
	window := M.windows[node]
	if len(window) < M.windowSize {
		window = append(window, rtt)
		M.windows[node] = window
		return rtt
	}
	window = append(window[1:], rtt)
	M.windows[node] = window
	samples := make([]time.Duration, M.windowSize)
	copy(samples, window)
	slices.Sort(samples)
	i := int(float64(len(samples)) * M.p)
	return samples[i]
}
