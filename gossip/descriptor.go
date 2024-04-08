package gossip

import "b4/shared"

type Descriptor struct {
	shared.Node
	Timestamp int
}

func NewDescriptor(node shared.Node, timestamp int) Descriptor {
	return Descriptor{Node: node, Timestamp: timestamp}
}
