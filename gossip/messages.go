package sampling

import (
	"b4/shared"
	"b4/vivaldi"
)

type MessageType uint8

const (
	Reply   MessageType = iota
	Request MessageType = iota
)

type PViewMessage struct {
	Type     MessageType
	Capacity int
	View     []Descriptor
	Coords   []RemoteCoord
}

type RemoteCoord struct {
	Owner shared.Node
	Coord vivaldi.Coord
	Age   int64
}
