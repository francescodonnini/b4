package gossip

import (
	"b4/shared"
	"b4/vivaldi"
	"time"
)

type RemoteCoord struct {
	Owner shared.Node
	Coord vivaldi.Coord
	Age   time.Time
}

func NewRemoteCoord(owner shared.Node, coord vivaldi.Coord, age time.Time) RemoteCoord {
	return RemoteCoord{Owner: owner, Coord: coord, Age: age}
}
