package gossip

import (
	"b4/shared"
	"b4/vivaldi"
)

type RemoteCoord struct {
	Owner shared.Node
	Coord vivaldi.Coord
	Age   int64
}

func NewRemoteCoord(owner shared.Node, coord vivaldi.Coord, age int64) RemoteCoord {
	return RemoteCoord{Owner: owner, Coord: coord, Age: age}
}
