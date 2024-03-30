package gossip

import (
	"b4/shared"
	"b4/vivaldi"
)

type Store interface {
	Peers() []shared.Node
	Read(node shared.Node) (vivaldi.Coord, bool)
	Save(coord RemoteCoord)
}

type StoreImpl struct {
	coords map[shared.Node]RemoteCoord
}

func NewStoreImpl() Store {
	return &StoreImpl{coords: make(map[shared.Node]RemoteCoord)}
}

func (s *StoreImpl) Peers() []shared.Node {
	peers := make([]shared.Node, 0)
	for k, _ := range s.coords {
		peers = append(peers, k)
	}
	return peers
}

func (s *StoreImpl) Read(node shared.Node) (vivaldi.Coord, bool) {
	c, ok := s.coords[node]
	if !ok {
		return vivaldi.Coord{}, ok
	}
	return c.Coord, ok
}

func (s *StoreImpl) Save(coord RemoteCoord) {
	c, ok := s.coords[coord.Owner]
	if !ok || c.Age < coord.Age {
		s.coords[c.Owner] = coord
	}
}
