package sampling

import (
	"b4/shared"
	"b4/vivaldi"
	"sync"
	"unsafe"
)

type Value struct {
	RemoteCoord
	Counter int
}

func NewValue(coord RemoteCoord, counter int) *Value {
	return &Value{coord, counter}
}

func (v *Value) Dec() {
	v.Counter--
}

type Spreader struct {
	mu    *sync.RWMutex
	cache map[shared.Node]*Value
	bus   *shared.EventBus
	maxN  int
}

func NewSpreader(bus *shared.EventBus) *Spreader {
	return &Spreader{
		bus:   bus,
		cache: make(map[shared.Node]*Value),
		mu:    &sync.RWMutex{},
		maxN:  6,
	}
}

func (s *Spreader) Select(maxNumOfBytes uintptr) (map[shared.Node]vivaldi.Coord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	values := make(map[shared.Node]vivaldi.Coord)
	var numOfBytes uintptr = 0
	for ip, v := range s.cache {
		if v.Counter == 0 {
			delete(s.cache, ip)
			continue
		}
		step := unsafe.Sizeof(v.Coord) + unsafe.Sizeof(ip)
		if numOfBytes+step > maxNumOfBytes {
			break
		}
		v.Dec()
		values[ip] = v.Coord
		numOfBytes += step
	}
	return values, nil
}

func (s *Spreader) Spread(coord RemoteCoord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.cache[coord.Owner]
	if !ok {
		s.cache[coord.Owner] = NewValue(coord, s.maxN)
		s.bus.Publish(shared.Event{
			Topic:   "coord/store",
			Content: coord,
		})
	}
}
