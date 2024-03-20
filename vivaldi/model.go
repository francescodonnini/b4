package vivaldi

import (
	"b4/shared"
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
	sampler    Sampler
}

func DefaultModel() Model {
	return NewModel(0.25, 0.25, NewMPFilter(4, 0.25))
}

func NewModel(cc, ce float64, sampler Sampler) Model {
	rand.NewSource(time.Now().Unix())
	return &ModelImpl{
		cc: cc,
		ce: ce,
		coord: Coord{
			Point:  NewVec3d(rand.Float64(), rand.Float64(), rand.Float64()),
			Height: 0,
		},
		localError: rand.Float64(),
		mu:         &sync.RWMutex{},
		sampler:    sampler,
	}
}

func (m ModelImpl) Update(rtt time.Duration, coord Coord, remoteError float64, node shared.Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	w := m.localError / (m.localError + remoteError)
	diff := m.coord.Sub(coord)
	dist := diff.Magnitude()
	// filter raw RTTs to remove the outliers. See:
	// http://nrs.harvard.edu/urn-3:HUL.InstRepos:25686820
	sample := m.sampler.Update(node, rtt)
	rttSeconds := sample.Seconds()
	relativeError := math.Abs(dist-rttSeconds) / rttSeconds
	// localError should be between 0 and 1. See:
	// http://nrs.harvard.edu/urn-3:HUL.InstRepos:25686820
	m.localError = math.Min(relativeError*m.ce*w+m.localError*(1-m.ce*w), 1)
	d := m.cc * w
	shift := diff.Unit().Scale((rttSeconds - dist) * d)
	m.coord = m.coord.Add(shift)
}

func (m ModelImpl) GetCoord() (Coord, float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.coord, m.localError
}
