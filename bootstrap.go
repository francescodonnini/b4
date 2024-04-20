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
	eventbus "github.com/francescodonnini/pubsub"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

func initialize(bus *eventbus.EventBus, store gossip.Store) {
	rand.NewSource(time.Now().Unix())
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
	filter := newFilter(settings)
	model := newVivaldiModel(bus, filter, settings)
	spreader, membership := newMembershipProtocol(peers, bus)
	energy := vivaldi.NewEnergySlidingWindow(16, 0.001, bus)
	// Si mette in ascolto delle coordinate di sistema per aggiornare la coordinata di applicazione.
	energyLis := bus.Subscribe("coord/sys")
	go func() {
		for e := range energyLis {
			energy.Update(e.Content.(vivaldi.Coord))
		}
	}()
	// Si mette in ascolto delle coordinate di applicazione per memorizzarle in Store.
	appLis := bus.Subscribe("coord/app")
	go func() {
		for e := range appLis {
			pair := e.Content.(shared.Pair[vivaldi.Coord, time.Time])
			spreader.Spread(gossip.NewRemoteCoord(id, pair.First, pair.Second))
		}
	}()
	// Ogni volta che si ricevono nuove coordinate vengono passate all'unità di gossiping per essere diffuse.
	gossipLis := bus.Subscribe("coord/update")
	go func() {
		for e := range gossipLis {
			updates := e.Content.([]gossip.RemoteCoord)
			spreader.Spread(updates...)
		}
	}()
	// Si mette in ascolto delle coordinate ricevute da altri nodi per memorizzarle nello store.
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
	go startVivaldiClient(membership, model, settings)
	go startGrpcServices(model, repl.NewShell(id, store), lis)
	go startUdpServer(context.Background(), id, filter, membership, bus)
	go startUdpClient(membership, settings)
	exitLis := bus.Subscribe("b4/exit")
	for range exitLis {
		os.Exit(0)
	}
}

func bootstrap(id shared.Node, discovery discv.Client) []shared.Node {
	// Si chiede la lista dei peer al registry fino a quando non se ne ottengono almeno 10
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
	timeout := settings.GetIntOrDefault("GOSSIP_TIMEOUT", 3000)
	client := vivaldi_grpc.NewClient(sampl, model, shared.NewDialer())
	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	for range ticker.C {
		client.Update()
	}
}

func startUdpClient(protocol gossip.Protocol, settings shared.Settings) {
	timeout := settings.GetIntOrDefault("GOSSIP_TIMEOUT", 3000)
	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	for range ticker.C {
		protocol.OnTimeout()
	}
}

func startHeartBeatClient(beat *discv.HeartBeatClient, settings shared.Settings) {
	timeout := settings.GetIntOrDefault("HEARTBEAT_TIMEOUT", 5000)
	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	for range ticker.C {
		beat.Beat()
	}
}

func newFilter(settings shared.Settings) shared.Filter {
	typ := settings.GetStringOrDefault("FILTER", "")
	if typ == "RAW" {
		return shared.NewRawFilter()
	}
	winSz := settings.GetIntOrDefault("MPF_WINDOW_SIZE", 16)
	p := settings.GetFloatOrDefault("MPF_PERCENTILE", 0.25)
	return shared.NewMPFilter(winSz, p)
}

func newVivaldiModel(bus *eventbus.EventBus, filter shared.Filter, settings shared.Settings) vivaldi.Model {
	cc := settings.GetFloatOrDefault("CC", 0.25)
	ce := settings.GetFloatOrDefault("CE", 0.25)
	dim := settings.GetIntOrDefault("DIMENSIONALITY", 3)
	return vivaldi.NewModel(cc, ce, dim, filter, bus)
}

func newMembershipProtocol(peers []shared.Node, bus *eventbus.EventBus) (*gossip.Spreader, *gossip.Impl) {
	settings := shared.NewSettings()
	maxN := settings.GetIntOrDefault("FEEDBACK_COUNTER", 6)
	id := shared.GetAddress()
	capacity := settings.GetIntOrDefault("CAPACITY", 8)
	spreader := gossip.NewSpreader(bus, maxN)
	return spreader, gossip.NewProtocol(id, capacity, peers, gossip.NewClient(spreader, id))
}
