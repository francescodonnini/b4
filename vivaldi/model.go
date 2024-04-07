package vivaldi

import (
	"b4/shared"
	"log"
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
	sampler    Filter
	bus        *shared.EventBus
}

func DefaultModel(bus *shared.EventBus) Model {
	return NewModel(0.25, 0.25, NewMPFilter(4, 0.25), bus)
}

func NewModel(cc, ce float64, sampler Filter, bus *shared.EventBus) Model {
	rand.NewSource(time.Now().Unix())
	return &ModelImpl{
		cc: cc,
		ce: ce,
		coord: Coord{
			Point: NewRandomUnit(),
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
	m.bus.Publish(shared.Event{
		Topic:   "coord/sys",
		Content: m.coord,
	})
	log.Printf("error %s %f\n", node.Address(), math.Abs(rtt.Seconds()-distance(m.coord, coord)))
}

func distance(x, y Coord) float64 {
	return x.Sub(y).Magnitude()
}

func (m *ModelImpl) GetCoord() (Coord, float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.coord, m.localError
}
