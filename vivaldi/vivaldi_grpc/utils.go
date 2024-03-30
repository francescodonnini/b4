package vivaldi_grpc

import (
	"b4/vivaldi"
	"b4/vivaldi/vivaldi_grpc/vivaldi_pb"
)

func (c GrpcClient) proto2model(coord *vivaldi_pb.Coord) (vivaldi.Coord, float64) {
	return vivaldi.Coord{
		Point: vivaldi.NewVec3d(
			coord.Point[0],
			coord.Point[1],
			coord.Point[2]),
	}, coord.Error
}
