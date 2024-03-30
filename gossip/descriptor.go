package gossip

import "b4/shared"

type Descriptor struct {
	shared.Node
	Timestamp int64
}

func NewDescriptor(node shared.Node, timestamp int64) Descriptor {
	return Descriptor{Node: node, Timestamp: timestamp}
}
