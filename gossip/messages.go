package gossip

import (
	"b4/shared"
	"time"
)

type MessageType uint8

const (
	Reply   MessageType = iota
	Request MessageType = iota
)

type PViewMessage struct {
	Type      MessageType
	Capacity  int
	View      []Descriptor
	Coords    []RemoteCoord
	Timestamp time.Time
	Srv       shared.Node
}

func NewReply(capacity int, view []Descriptor, coords []RemoteCoord, timestamp time.Time, srv shared.Node) PViewMessage {
	return PViewMessage{Type: Reply, Capacity: capacity, View: view, Coords: coords, Timestamp: timestamp, Srv: srv}
}

func NewRequest(capacity int, view []Descriptor, coords []RemoteCoord, timestamp time.Time, srv shared.Node) PViewMessage {
	return PViewMessage{Type: Request, Capacity: capacity, View: view, Coords: coords, Timestamp: timestamp, Srv: srv}
}
