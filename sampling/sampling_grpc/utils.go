package sampling_grpc

import (
	"b4/sampling"
	"b4/sampling/sampling_grpc/sampling_pb"
	"b4/shared"
)

func model2proto(view sampling.PView, source shared.Node) *sampling_pb.View {
	descriptors := make([]*sampling_pb.Descriptor, 0)
	for _, desc := range view.Descriptors() {
		descriptors = append(descriptors, &sampling_pb.Descriptor{
			Ip:   desc.Ip,
			Port: int32(desc.Port),
			Age:  desc.Age,
		})
	}
	return &sampling_pb.View{
		Descriptors: descriptors,
		C:           int32(view.Capacity()),
		SourceIp:    source.Ip,
		SourcePort:  int32(source.Port),
	}
}

func proto2model(descriptors []*sampling_pb.Descriptor, c int) sampling.PView {
	view := make([]sampling.Descriptor, 0)
	for _, desc := range descriptors {
		node := shared.Node{
			Ip:   desc.Ip,
			Port: int(desc.Port),
		}
		view = append(view, sampling.NewDescriptor(node, desc.Age))
	}
	return sampling.NewView(c, view)
}
