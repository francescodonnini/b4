package gossip

import (
	"b4/shared"
	"b4/vivaldi"
	"log"
	"sync"
)

type Store interface {
	Peers() []shared.Node
	Items() []RemoteCoord
	Read(node shared.Node) (vivaldi.Coord, bool)
	Save(coord RemoteCoord)
}

type StoreImpl struct {
	mu     *sync.RWMutex
	coords map[shared.Node]RemoteCoord
}

func NewStoreImpl() Store {
	return &StoreImpl{
		mu:     &sync.RWMutex{},
		coords: make(map[shared.Node]RemoteCoord)}
}

func (s *StoreImpl) Items() []RemoteCoord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]RemoteCoord, 0)
	for _, v := range s.coords {
		items = append(items, v)
	}
	return items
}

func (s *StoreImpl) Peers() []shared.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	peers := make([]shared.Node, 0)
	for k, _ := range s.coords {
		peers = append(peers, k)
	}
	return peers
}

func (s *StoreImpl) Read(node shared.Node) (vivaldi.Coord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.coords[node]
	if !ok {
		return vivaldi.Coord{}, ok
	}
	return c.Coord, ok
}

func (s *StoreImpl) Save(coord RemoteCoord) {
	s.mu.Lock()
	s.mu.Unlock()
	c, ok := s.coords[coord.Owner]
	if !ok || c.Age < coord.Age {
		s.coords[coord.Owner] = coord
		log.Printf("coord/store: %v\n", coord)
	}
}
