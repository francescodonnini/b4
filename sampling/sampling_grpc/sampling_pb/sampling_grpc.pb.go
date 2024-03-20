// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: sampling/sampling_grpc/sampling.proto

package sampling_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SamplingClient is the client API for Sampling service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SamplingClient interface {
	Exchange(ctx context.Context, in *View, opts ...grpc.CallOption) (*View, error)
}

type samplingClient struct {
	cc grpc.ClientConnInterface
}

func NewSamplingClient(cc grpc.ClientConnInterface) SamplingClient {
	return &samplingClient{cc}
}

func (c *samplingClient) Exchange(ctx context.Context, in *View, opts ...grpc.CallOption) (*View, error) {
	out := new(View)
	err := c.cc.Invoke(ctx, "/Sampling/Exchange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SamplingServer is the server API for Sampling service.
// All implementations must embed UnimplementedSamplingServer
// for forward compatibility
type SamplingServer interface {
	Exchange(context.Context, *View) (*View, error)
	mustEmbedUnimplementedSamplingServer()
}

// UnimplementedSamplingServer must be embedded to have forward compatible implementations.
type UnimplementedSamplingServer struct {
}

func (UnimplementedSamplingServer) Exchange(context.Context, *View) (*View, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exchange not implemented")
}
func (UnimplementedSamplingServer) mustEmbedUnimplementedSamplingServer() {}

// UnsafeSamplingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SamplingServer will
// result in compilation errors.
type UnsafeSamplingServer interface {
	mustEmbedUnimplementedSamplingServer()
}

func RegisterSamplingServer(s grpc.ServiceRegistrar, srv SamplingServer) {
	s.RegisterService(&Sampling_ServiceDesc, srv)
}

func _Sampling_Exchange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(View)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SamplingServer).Exchange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Sampling/Exchange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SamplingServer).Exchange(ctx, req.(*View))
	}
	return interceptor(ctx, in, info, handler)
}

// Sampling_ServiceDesc is the grpc.ServiceDesc for Sampling service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sampling_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Sampling",
	HandlerType: (*SamplingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Exchange",
			Handler:    _Sampling_Exchange_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sampling/sampling_grpc/sampling.proto",
}
