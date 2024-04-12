package gossip

import (
	"b4/shared"
	"time"
)

type MessageType uint8

const (
	Reply   = iota
	Request = iota
)

type PViewMessage struct {
	Capacity int
	View     []Descriptor
}

type Message struct {
	PViewMessage
	Coords    []RemoteCoord
	Source    shared.Node
	Timestamp time.Time
	Type      MessageType
}

func NewReply(view PViewMessage, coords []RemoteCoord, timestamp time.Time, source shared.Node) Message {
	return Message{
		PViewMessage: view,
		Coords:       coords,
		Source:       source,
		Timestamp:    timestamp,
		Type:         Reply,
	}
}

func NewRequest(view PViewMessage, coords []RemoteCoord, timestamp time.Time, source shared.Node) Message {
	return Message{
		PViewMessage: view,
		Coords:       coords,
		Source:       source,
		Timestamp:    timestamp,
		Type:         Request,
	}
}
