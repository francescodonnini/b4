package main

import (
	"b4/gossip"
	"b4/shared"
	"b4/vivaldi"
	eventbus "github.com/francescodonnini/pubsub"
	"log"
	"time"
)

func main() {
	config := shared.NewSettings()
	retention := config.GetIntOrDefault("RETENTION", 3)
	bus := eventbus.NewEventBus()
	store := gossip.NewRetentionPolicy(gossip.NewInMemoryStore(), time.Duration(retention)*time.Minute)
	go initialize(bus, store)
	coords := bus.Subscribe("coord/app")
	for e := range coords {
		p := e.Content.(shared.Pair[vivaldi.Coord, time.Time])
		log.Printf("coord/app: %v\n", p.First)
	}
}
