package vivaldi_grpc

import (
	"b4/vivaldi"
	"b4/vivaldi/vivaldi_grpc/vivaldi_pb"
)

func (c GrpcClient) proto2model(coord *vivaldi_pb.Coord) (vivaldi.Coord, float64) {
	return vivaldi.Coord{
		Point: coord.Point,
	}, coord.Error
}
