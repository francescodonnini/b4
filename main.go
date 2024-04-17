package main

import (
	"b4/gossip"
	"b4/shared"
	"b4/vivaldi"
	eventbus "github.com/francescodonnini/pubsub"
	"time"
)

func main() {
	bus := eventbus.NewEventBus()
	store := gossip.NewInMemoryStore()
	go initialize(bus, store)
	coords := bus.Subscribe("coord/app")
	for e := range coords {
		_ = e.Content.(shared.Pair[vivaldi.Coord, time.Time])
	}
}
