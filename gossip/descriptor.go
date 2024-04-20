package gossip

import "b4/shared"

type Descriptor struct {
	shared.Node
	Age int
}

func NewDescriptor(node shared.Node, age int) Descriptor {
	return Descriptor{Node: node, Age: age}
}
