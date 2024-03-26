package shared

import (
	"b4/gossip"
	"b4/vivaldi"
)

type Store interface {
	Read(node Node) (vivaldi.Coord, bool)
	Save(coord gossip.RemoteCoord)
}

type StoreImpl struct {
	coords map[Node]gossip.RemoteCoord
}

func NewStoreImpl() Store {
	return &StoreImpl{coords: make(map[Node]gossip.RemoteCoord)}
}

func (s *StoreImpl) Read(node Node) (vivaldi.Coord, bool) {
	c, ok := s.coords[node]
	if !ok {
		return vivaldi.Coord{}, ok
	}
	return c.Coord, ok
}

func (s *StoreImpl) Save(coord gossip.RemoteCoord) {
	c, ok := s.coords[coord.Owner]
	if !ok || c.Age < coord.Age {
		s.coords[c.Owner] = coord
	}
}
