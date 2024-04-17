package vivaldi_grpc

import (
	"b4/shared"
	"b4/vivaldi"
	"b4/vivaldi/vivaldi_grpc/vivaldi_pb"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type GrpcClient struct {
	sampling shared.PeerSampling
	model    vivaldi.Model
	dialer   shared.Dialer
}

func NewClient(sampling shared.PeerSampling, model vivaldi.Model, dialer shared.Dialer) GrpcClient {
	return GrpcClient{
		sampling: sampling,
		model:    model,
		dialer:   dialer,
	}
}

func (c GrpcClient) Update() {
	p, ok := c.sampling.GetRandom()
	if !ok {
		return
	}
	client, err := c.connect(p)
	if err != nil {
		return
	}
	start := time.Now()
	res, err := client.GetCoord(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Printf("grpc call GetCoord to %s failed. Error: %s\n", p.Address(), err)
		return
	}
	rtt := time.Now().Sub(start)
	coord, remoteError := c.proto2model(res)
	c.model.Update(rtt, coord, remoteError, p)
}

func (c GrpcClient) connect(node shared.Node) (vivaldi_pb.VivaldiClient, error) {
	conn, err := c.dialer.Dial(node)
	if err != nil {
		return nil, err
	}
	return vivaldi_pb.NewVivaldiClient(conn), nil
}
