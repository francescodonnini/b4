package main

import (
	"b4/gossip"
	"b4/shared"
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
	peers := bus.Subscribe("coord/store")
	// Ascolta le coordinate che si ottengono dagli altri nodi tramite gossiping e ci calcola le distanze.
	id := shared.GetAddress()
	for e := range peers {
		other := e.Content.(gossip.RemoteCoord)
		self, ok := store.Read(id)
		if ok && self.Owner != other.Owner {
			x := self.Coord
			y := other.Coord
			dist := x.Sub(y).Magnitude()
			log.Printf("rtt between this and %s is %f\n", other.Owner.Ip, dist)
		}
	}
}
