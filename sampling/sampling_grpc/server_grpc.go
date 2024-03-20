package sampling_grpc

import (
	"b4/sampling"
	"b4/sampling/sampling_grpc/sampling_pb"
	"b4/shared"
	"context"
)

type GrpcServer struct {
	sampling_pb.UnimplementedSamplingServer
	id      shared.Node
	service sampling.PeerSamplingService
}

func NewServer(id shared.Node, service sampling.PeerSamplingService) sampling_pb.SamplingServer {
	return &GrpcServer{
		id:      id,
		service: service,
	}
}

func (g GrpcServer) Exchange(_ context.Context, in *sampling_pb.View) (*sampling_pb.View, error) {
	view := g.service.OnReceive(proto2model(in.Descriptors, int(in.C)))
	return model2proto(view, g.id), nil
}
