package discovery

import (
	"b4/shared"
)

type Discovery interface {
	GetNodes() []shared.Node
}

type Client interface {
	Join(node shared.Node) []shared.Node
	Exit(node shared.Node)
}

type Service struct {
	id     shared.Node
	client Client
}

func NewDiscoveryService(endpoint, id shared.Node) Discovery {
	client := NewDiscoveryClient(endpoint)
	return Service{
		id:     id,
		client: client,
	}
}

func (d Service) GetNodes() []shared.Node {
	return d.client.Join(d.id)
}
