package sampling_grpc

import (
	"b4/sampling"
	"b4/sampling/sampling_grpc/sampling_pb"
	"b4/shared"
	"context"
	"log"
)

type GrpcClient struct {
	id     shared.Node
	dialer shared.Dialer
}

func NewClient(id shared.Node, dialer shared.Dialer) *GrpcClient {
	return &GrpcClient{
		id:     id,
		dialer: dialer,
	}
}

func (g *GrpcClient) Exchange(view sampling.PView, dest shared.Node) (sampling.PView, error) {
	client, err := g.connect(dest)
	if err != nil {
		return nil, err
	}
	other, err := client.Exchange(context.Background(), model2proto(view, g.id))
	if err != nil {
		log.Printf("grpc call to %s Exchange failed. Error: %s\n", dest.Address(), err)
		return nil, err
	}
	return proto2model(other.Descriptors, int(other.C)), nil
}

func (g *GrpcClient) connect(node shared.Node) (sampling_pb.SamplingClient, error) {
	conn, err := g.dialer.Dial(node)
	if err != nil {
		return nil, err
	}
	return sampling_pb.NewSamplingClient(conn), nil
}
