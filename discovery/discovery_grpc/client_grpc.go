package discovery_grpc

import (
	"b4/discovery/discovery_grpc/discovery_pb"
	"b4/shared"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type GrpcClient struct {
	client   discovery_pb.DiscoveryClient
	endpoint shared.Node
}

func NewDiscoveryClient(endpoint shared.Node) GrpcClient {
	return GrpcClient{
		endpoint: endpoint,
		client:   nil,
	}
}

func (g GrpcClient) Join(node shared.Node) []shared.Node {
	client, err := g.connect(g.endpoint)
	if err != nil {
		return make([]shared.Node, 0)
	}
	res, err := client.Join(context.Background(), &discovery_pb.Node{
		Ip:   node.Ip,
		Port: int32(node.Port),
	})
	if err != nil {
		return make([]shared.Node, 0)
	}
	peers := make([]shared.Node, 0)
	for _, p := range res.Peers {
		peers = append(peers, shared.Node{
			Ip:   p.Ip,
			Port: int(p.Port),
		})
	}
	return peers
}

func (g GrpcClient) Exit(node shared.Node) {
	client, err := g.connect(g.endpoint)
	if err != nil {
		return
	}
	_, err = client.Exit(context.Background(), &discovery_pb.Node{
		Ip:   node.Ip,
		Port: int32(node.Port),
	})
}

func (g GrpcClient) connect(node shared.Node) (discovery_pb.DiscoveryClient, error) {
	if g.client != nil {
		return g.client, nil
	}
	conn, err := grpc.Dial(node.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Cannot connect to %s. Error: %s\n", node.Address(), err)
		return nil, err
	}
	client := discovery_pb.NewDiscoveryClient(conn)
	g.client = client
	return client, nil
}
