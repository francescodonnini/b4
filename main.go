package main

import (
	discv "b4/discovery"
	"b4/gossip"
	"b4/repl"
	"b4/repl/shell_grpc"
	"b4/repl/shell_grpc/shell_pb"
	"b4/shared"
	"b4/vivaldi"
	"b4/vivaldi/vivaldi_grpc"
	"b4/vivaldi/vivaldi_grpc/vivaldi_pb"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	if enabled, err := strconv.ParseBool(os.Getenv("LOGGING_ENABLED")); err != nil || enabled == false {
		log.SetOutput(io.Discard)
	}
	ip := shared.GetIp()
	id := shared.Node{Ip: ip.String(), Port: 5050}
	// TODO: Leggere IP registry da file di configurazione.
	endpoint := shared.Node{Ip: "10.0.0.253", Port: 5050}
	discovery := discv.NewDiscoveryService(endpoint, id)
	_ = discovery.GetNodes()
	go startHeartBeatClient(discv.NewHeartBeatClient(shared.Node{Ip: "10.0.0.253", Port: 5050}))
	peers := bootstrap(id, discovery)
	bus := shared.NewEventBus()
	spreader := gossip.NewSpreader(bus, 6)
	membership := gossip.NewProtocol(id, 4, peers, gossip.NewClient(spreader))
	model := vivaldi.DefaultModel(bus)
	energy := vivaldi.NewEnergySlidingWindow(16, 0.001, bus)
	energyLis := bus.Subscribe("coord/sys")
	go func() {
		for e := range energyLis {
			energy.Update(e.Content.(vivaldi.Coord))
		}
	}()
	appLis := bus.Subscribe("coord/app")
	go func() {
		for e := range appLis {
			pair := e.Content.(shared.Pair[vivaldi.Coord, int64])
			log.Printf("coord/app: %v\n", pair.First)
			spreader.Spread(gossip.NewRemoteCoord(id, pair.First, pair.Second))
		}
	}()
	gossipLis := bus.Subscribe("coord/update")
	go func() {
		for e := range gossipLis {
			updates := e.Content.([]gossip.RemoteCoord)
			spreader.Spread(updates...)
		}
	}()
	store := gossip.NewStoreImpl()
	storeLis := bus.Subscribe("coord/store")
	go func() {
		for e := range storeLis {
			coord := e.Content.(gossip.RemoteCoord)
			store.Save(coord)
		}
	}()
	lis, err := net.Listen("tcp", id.Address())
	if err != nil {
		log.Fatalf("Failed to listen: %s\n", err)
	}
	go startVivaldiClient(membership, model)
	go startGrpcServices(model, repl.NewShell(id, store), lis)
	go startUdpServer(context.Background(), id, membership, bus)
	go startUdpClient(membership)
	select {}
}

func bootstrap(id shared.Node, discovery discv.Discovery) []shared.Node {
	peers := discovery.GetNodes()
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		if len(peers) >= 10 {
			break
		}
		peers = discovery.GetNodes()
	}
	peers = shared.RemoveIf(peers, func(node shared.Node) bool {
		return node == id
	})
	return peers
}

func startUdpServer(ctx context.Context, id shared.Node, membership gossip.Protocol, bus *shared.EventBus) {
	srv := gossip.NewUdpServer(id, membership, bus)
	srv.Serve(ctx)
}

func startGrpcServices(model vivaldi.Model, shell repl.Shell, lis net.Listener) {
	srv := grpc.NewServer()
	vivaldi_pb.RegisterVivaldiServer(srv, vivaldi_grpc.NewServer(model))
	shell_pb.RegisterShellServer(srv, shell_grpc.NewGrpcServer(shell))
	err := srv.Serve(lis)
	if err != nil {
		log.Fatalf("Cannot serve: %s\n", err)
	}
}

func startVivaldiClient(sampl shared.PeerSampling, model vivaldi.Model) {
	client := vivaldi_grpc.NewClient(sampl, model, shared.NewDialer())
	for range time.Tick(3 * time.Second) {
		client.Update()
	}
}

func startUdpClient(protocol gossip.Protocol) {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		protocol.OnTimeout()
	}
}

func startHeartBeatClient(beat *discv.HeartBeatClient) {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		beat.Beat()
	}
}
