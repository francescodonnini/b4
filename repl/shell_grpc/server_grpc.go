package shell_grpc

import (
	"b4/repl"
	"b4/repl/shell_grpc/shell_pb"
	"context"
)

type GrpcServer struct {
	shell_pb.UnimplementedShellServer
	shell repl.Shell
}

func NewGrpcServer(shell repl.Shell) *GrpcServer {
	return &GrpcServer{shell: shell}
}

func (g *GrpcServer) Execute(_ context.Context, in *shell_pb.Line) (*shell_pb.Output, error) {
	res, err := g.shell.Execute(in.Text)
	return &shell_pb.Output{Value: string(res[:])}, err
}
