package gossip

import (
	"b4/shared"
	"sync"
)

type Store interface {
	Peers() []shared.Node
	Items() []RemoteCoord
	Read(node shared.Node) (RemoteCoord, bool)
	Save(coord RemoteCoord)
}

type InMemoryStore struct {
	mu     *sync.RWMutex
	coords map[shared.Node]RemoteCoord
}

func NewInMemoryStore() Store {
	return &InMemoryStore{
		mu:     &sync.RWMutex{},
		coords: make(map[shared.Node]RemoteCoord)}
}

func (s *InMemoryStore) Items() []RemoteCoord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]RemoteCoord, 0)
	for _, v := range s.coords {
		items = append(items, v)
	}
	return items
}

func (s *InMemoryStore) Peers() []shared.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	peers := make([]shared.Node, 0)
	for k, _ := range s.coords {
		peers = append(peers, k)
	}
	return peers
}

func (s *InMemoryStore) Read(node shared.Node) (RemoteCoord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.coords[node]
	if !ok {
		return RemoteCoord{}, ok
	}
	return c, ok
}

func (s *InMemoryStore) Save(coord RemoteCoord) {
	s.mu.Lock()
	s.mu.Unlock()
	c, ok := s.coords[coord.Owner]
	if !ok || c.Age.Before(coord.Age) {
		s.coords[coord.Owner] = coord
	}
}
