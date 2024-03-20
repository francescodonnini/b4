package main

import (
	"b4/discovery"
	"b4/sampling"
	"b4/sampling/sampling_grpc"
	"b4/sampling/sampling_grpc/sampling_pb"
	"b4/shared"
	"b4/vivaldi"
	"b4/vivaldi/vivaldi_grpc"
	"b4/vivaldi/vivaldi_grpc/vivaldi_pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	ip := shared.GetIp()
	id := shared.Node{
		Ip:   ip.String(),
		Port: 5050,
	}
	endpoint := shared.Node{
		Ip:   "10.0.0.253",
		Port: 5050,
	}
	disc := discovery.NewDiscoveryService(endpoint, id)
	peers := disc.GetNodes()
	for range time.Tick(3 * time.Second) {
		if len(peers) >= 7 {
			break
		}
		peers = disc.GetNodes()
	}

	lis, err := net.Listen("tcp", id.Address())
	if err != nil {
		log.Fatalf("Failed to listen: %s\n", err)
	}
	dialer := shared.NewDialer()
	samplingClient := sampling_grpc.NewClient(id, dialer)
	sampl := sampling.NewSamplingService(id, 3, peers, samplingClient)
	model := vivaldi.DefaultModel()

	server := grpc.NewServer()
	registerServices(server, model, id, sampl)
	go startServer(server, lis)

	done := make(chan bool)
	go startSamplingService(sampl)
	go startVivaldiService(vivaldi_grpc.NewClient(sampl, model, dialer))
	<-done
}

func startVivaldiService(client vivaldi.Client) {
	for range time.Tick(3 * time.Second) {
		client.Update()
	}
}

func registerServices(s *grpc.Server, model vivaldi.Model, id shared.Node, sampl sampling.PeerSamplingService) {
	sampling_pb.RegisterSamplingServer(s, sampling_grpc.NewServer(id, sampl))
	vivaldi_pb.RegisterVivaldiServer(s, vivaldi_grpc.NewServer(model))
}

func startServer(s *grpc.Server, lis net.Listener) {
	err := s.Serve(lis)
	if err != nil {
		log.Fatalf("Cannot serve: %s\n", err)
	}
}

func startSamplingService(service sampling.PeerSamplingService) {
	for range time.Tick(time.Second * 5) {
		service.OnTimeout()
	}
}
