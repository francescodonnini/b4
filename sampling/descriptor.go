package sampling

import "b4/shared"

type Descriptor struct {
	shared.Node
	Age int64
}

func FreshDescriptor(node shared.Node) Descriptor {
	return NewDescriptor(node, 0)
}

func NewDescriptor(node shared.Node, age int64) Descriptor {
	return Descriptor{
		Node: node,
		Age:  age,
	}
}
