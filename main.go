package main

import (
	discv "b4/discovery"
	"b4/gossip"
	"b4/logging"
	"b4/repl"
	"b4/repl/shell_grpc"
	"b4/repl/shell_grpc/shell_pb"
	"b4/shared"
	"b4/vivaldi"
	"b4/vivaldi/vivaldi_grpc"
	"b4/vivaldi/vivaldi_grpc/vivaldi_pb"
	"context"
	eventbus "github.com/francescodonnini/pubsub"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	settings := shared.NewSettings()
	if enabled, ok := settings.GetBool("LOGGING_ENABLED"); !ok || enabled == false {
		log.SetOutput(io.Discard)
	}
	id := shared.GetAddress()
	endpoint := shared.GetRegistryAddress()
	discClient := discv.NewDiscoveryClient(endpoint)
	discClient.Join(id)
	go startHeartBeatClient(discv.NewHeartBeatClient(endpoint), settings)
	peers := bootstrap(id, discClient)
	bus := eventbus.NewEventBus()
	filter := newFilter(settings)
	model := newVivaldiModel(bus, filter, settings)
	spreader, membership := newMembershipProtocol(peers, bus)
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
			pair := e.Content.(shared.Pair[vivaldi.Coord, time.Time])
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
			logging.GetInstance("store").Printf("%v\n", coord)
		}
	}()
	lis, err := net.Listen("tcp", id.Address())
	if err != nil {
		log.Fatalf("Failed to listen: %s\n", err)
	}
	go startVivaldiClient(membership, model, settings)
	go startGrpcServices(model, repl.NewShell(id, store), lis)
	go startUdpServer(context.Background(), id, filter, membership, bus)
	go startUdpClient(membership, settings)
	time.Sleep(30 * time.Minute)
	os.Exit(0)
}

func bootstrap(id shared.Node, discovery discv.Client) []shared.Node {
	peers := discovery.Join(id)
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		if len(peers) >= 10 {
			break
		}
		peers = discovery.Join(id)
	}
	peers = shared.RemoveIf(peers, func(node shared.Node) bool {
		return node == id
	})
	return peers
}

func startUdpServer(ctx context.Context, id shared.Node, filter shared.Filter, membership gossip.Protocol, bus *eventbus.EventBus) {
	srv := gossip.NewUdpServer(id, membership, filter, bus)
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

func startVivaldiClient(sampl shared.PeerSampling, model vivaldi.Model, settings shared.Settings) {
	timeout, ok := settings.GetInt("GOSSIP_TIMEOUT")
	if !ok {
		timeout = 3000
	}
	client := vivaldi_grpc.NewClient(sampl, model, shared.NewDialer())
	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	for range ticker.C {
		client.Update()
	}
}

func startUdpClient(protocol gossip.Protocol, settings shared.Settings) {
	timeout, ok := settings.GetInt("GOSSIP_TIMEOUT")
	if !ok {
		timeout = 3000
	}
	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	for range ticker.C {
		protocol.OnTimeout()
	}
}

func startHeartBeatClient(beat *discv.HeartBeatClient, settings shared.Settings) {
	timeout, ok := settings.GetInt("HEARTBEAT_TIMEOUT")
	if !ok {
		timeout = 5000
	}
	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	for range ticker.C {
		beat.Beat()
	}
}

func newFilter(settings shared.Settings) shared.Filter {
	typ, ok := settings.GetString("FILTER")
	if !ok {
		typ = "MP"
	}
	if typ == "RAW" {
		log.Println("raw")
		return shared.NewRawFilter()
	}
	winSz, ok := settings.GetInt("MPF_WINDOW_SIZE")
	if !ok {
		winSz = 16
	}
	p, ok := settings.GetFloat("MPF_PERCENTILE")
	if !ok {
		p = 0.25
	}
	return shared.NewMPFilter(winSz, p)
}

func newVivaldiModel(bus *eventbus.EventBus, filter shared.Filter, settings shared.Settings) vivaldi.Model {
	cc, ok := settings.GetFloat("CC")
	if !ok {
		cc = 0.25
	}
	ce, ok := settings.GetFloat("CE")
	if !ok {
		ce = 0.25
	}
	dim, ok := settings.GetInt("DIMENSIONALITY")
	if !ok {
		dim = 3
	}
	return vivaldi.NewModel(cc, ce, dim, filter, bus)
}

func newMembershipProtocol(peers []shared.Node, bus *eventbus.EventBus) (*gossip.Spreader, *gossip.Impl) {
	settings := shared.NewSettings()
	maxN, ok := settings.GetInt("FEEDBACK_COUNTER")
	if !ok {
		maxN = 6
	}
	id := shared.GetAddress()
	capacity, ok := settings.GetInt("CAPACITY")
	if !ok {
		capacity = 6
	}
	spreader := gossip.NewSpreader(bus, maxN)
	return spreader, gossip.NewProtocol(id, capacity, peers, gossip.NewClient(spreader, id))
}
