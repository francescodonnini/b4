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
	retention := config.GetIntOrDefault("RETENTION", 15)
	bus := eventbus.NewEventBus()
	store := gossip.NewRetentionPolicy(gossip.NewInMemoryStore(), time.Duration(retention)*time.Minute)
	go initialize(bus, store)
	coords := bus.Subscribe("coord/app")
	peers := bus.Subscribe("coord/store")
	go func() {
		id := shared.GetAddress()
		for e := range peers {
			other := e.Content.(gossip.RemoteCoord)
			self, ok := store.Read(id)
			if ok {
				x := self.Coord
				y := other.Coord
				dist := x.Sub(y).Magnitude()
				log.Printf("rtt between this and %s is %f\n", other.Owner.Ip, dist)
			}
		}
	}()
	go func() {
		for e := range coords {
			p := e.Content.(shared.Pair[vivaldi.Coord, time.Time])
			log.Printf("coord/app: %v\n", p.First)
		}
	}()
	select {}
}
