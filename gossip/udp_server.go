package gossip

import (
	"b4/shared"
	"bytes"
	"context"
	"encoding/gob"
	event_bus "github.com/francescodonnini/pubsub"
	"log"
	"net"
	"time"
)

type UdpServer struct {
	id       shared.Node
	sampling Protocol
	filter   shared.Filter
	bus      *event_bus.EventBus
}

func NewUdpServer(id shared.Node, sampling Protocol, filter shared.Filter, bus *event_bus.EventBus) *UdpServer {
	return &UdpServer{id: id, sampling: sampling, bus: bus, filter: filter}
}

// Serve mette UdpServer in attesa di pacchetti UDP all'indirizzo specificato nel campo id.
// Il server si aspetta due tipi di messaggi che possono essere di richiesta (della vista parziale del nodo) e di
// risposta (a una richiesta inviata in precedenza). La gestione dei messaggi è delegata a sampling.
func (s *UdpServer) Serve(ctx context.Context) {
	address := s.id.Address()
	srv, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Fatalf("cannot listen to %s, error: %s\n", address, err)
	}
	go func() {
		go func() {
			<-ctx.Done()
			_ = srv.Close()
		}()
		buf := make([]byte, 65536)
		for {
			n, _, err := srv.ReadFrom(buf)
			if err != nil {
				log.Printf("cannot read from udp socket, error: %s\n", err)
				return
			}
			message, err := decodeMessage(buf[:n])
			if err != nil {
				continue
			}
			s.bus.Publish(event_bus.Event{Topic: "coord/update", Content: message.Coords})
			message.View = s.removeSelf(message.View)
			if message.Type == Reply {
				rtt := time.Now().Sub(message.Timestamp)
				s.filter.Update(message.Srv, rtt)
				s.sampling.OnReceiveReply(message)
			} else if message.Type == Request {
				s.sampling.OnReceiveRequest(message)
			}
		}
	}()
}

func (s *UdpServer) removeSelf(view []Descriptor) []Descriptor {
	return shared.RemoveIf(view, func(descriptor Descriptor) bool {
		return descriptor.Node == s.id
	})
}

func decodeMessage(payload []byte) (PViewMessage, error) {
	var message PViewMessage
	dec := gob.NewDecoder(bytes.NewReader(payload))
	err := dec.Decode(&message)
	if err != nil {
		log.Printf("Cannot decode request. Error: %s\n", err)
		return PViewMessage{}, err
	}
	switch message.Type {
	case Reply, Request:
		return message, nil
	default:
		log.Printf("Unknown message type= %d\n", message.Type)
	}
	return PViewMessage{}, err
}
