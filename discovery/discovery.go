package discovery

import (
	"b4/shared"
)

type Client interface {
	Join(node shared.Node) []shared.Node
	Exit(node shared.Node)
}

type Service struct {
	id     shared.Node
	client Client
}

func (d Service) GetNodes() []shared.Node {
	return d.client.Join(d.id)
}
