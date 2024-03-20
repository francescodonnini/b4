package discovery

import (
	"b4/shared"
	"fmt"
)

type MockDiscovery struct{}

func (m *MockDiscovery) GetNodes() []shared.Node {
	nodes := make([]shared.Node, 0)
	for i := 1; i <= 6; i++ {
		nodes = append(nodes, shared.Node{
			Ip:   fmt.Sprintf("10.0.0.%d", i),
			Port: 5050,
		})
	}
	return nodes
}
