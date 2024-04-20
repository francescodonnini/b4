package vivaldi

import (
	"b4/shared"
	eventbus "github.com/francescodonnini/pubsub"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Model interface {
	Update(rtt time.Duration, coord Coord, remoteError float64, node shared.Node)
	GetCoord() (Coord, float64)
}

type ModelImpl struct {
	cc         float64
	ce         float64
	coord      Coord
	localError float64
	mu         *sync.RWMutex
	sampler    shared.Filter
	bus        *eventbus.EventBus
}

func NewModel(cc, ce float64, n int, sampler shared.Filter, bus *eventbus.EventBus) Model {
	return &ModelImpl{
		cc: cc,
		ce: ce,
		coord: Coord{
			Point: NewRandomUnit(n),
		},
		localError: rand.Float64(),
		mu:         &sync.RWMutex{},
		sampler:    sampler,
		bus:        bus,
	}
}

func (m *ModelImpl) Update(rtt time.Duration, coord Coord, remoteError float64, node shared.Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	w := m.localError / (m.localError + remoteError)
	diff := m.coord.Sub(coord)
	dist := diff.Magnitude()
	// Filtra gli rtt per compensare i ritardi di rete. Vedere:
	// http://nrs.harvard.edu/urn-3:HUL.InstRepos:25686820
	sample := m.sampler.Update(node, rtt)
	rttSeconds := sample.Seconds()
	relativeError := math.Abs(dist-rttSeconds) / rttSeconds
	// localError deve essere compreso tra 0 e 1. Vedere:
	m.localError = math.Min(relativeError*m.ce*w+m.localError*(1-m.ce*w), 1)
	d := m.cc * w
	shift := diff.Unit().Scale((rttSeconds - dist) * d)
	m.coord = m.coord.Add(shift)
	m.bus.Publish(eventbus.Event{
		Topic:   "coord/sys",
		Content: m.coord,
	})
}

func distance(x, y Coord) float64 {
	return x.Sub(y).Magnitude()
}

func (m *ModelImpl) GetCoord() (Coord, float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.coord, m.localError
}
