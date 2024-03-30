package vivaldi_grpc

import (
	"b4/vivaldi"
	"b4/vivaldi/vivaldi_grpc/vivaldi_pb"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	vivaldi_pb.UnimplementedVivaldiServer
	model vivaldi.Model
}

func NewServer(model vivaldi.Model) vivaldi_pb.VivaldiServer {
	return &GrpcServer{model: model}
}

func (s GrpcServer) GetCoord(_ context.Context, _ *emptypb.Empty) (*vivaldi_pb.Coord, error) {
	coord, err := s.model.GetCoord()
	return &vivaldi_pb.Coord{
		Point: []float64{coord.Point[0], coord.Point[1], coord.Point[2]},
		Error: err,
	}, nil
}
