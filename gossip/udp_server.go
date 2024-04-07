package gossip

import (
	"b4/shared"
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"net"
	"strings"
)

type UdpServer struct {
	id       shared.Node
	sampling Protocol
	bus      *shared.EventBus
}

func NewUdpServer(id shared.Node, sampling Protocol, bus *shared.EventBus) *UdpServer {
	return &UdpServer{id: id, sampling: sampling, bus: bus}
}

// Serve mette UdpServer in attesa di pacchetti UDP all'indirizzo specificato nel campo id.
// Il server si aspetta due tipi di messaggi che possono essere di richiesta (della vista parziale del nodo) e di
// risposta (a una richiesta inviata in precedenza). La gestione dei messaggi è delegata a sampling.
func (s *UdpServer) Serve(ctx context.Context) {
	address := s.id.Address()
	srv, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Fatalf("Cannot listen to %s. Error: %s\n", address, err)
	}
	go func() {
		go func() {
			<-ctx.Done()
			_ = srv.Close()
		}()
		buf := make([]byte, 65536)
		for {
			n, addr, err := srv.ReadFrom(buf)
			if err != nil {
				log.Printf("cannot read from udp socket. error: %s\n", err)
				return
			}
			message, err := decodeMessage(buf[:n])
			if err != nil {
				continue
			}
			s.bus.Publish(shared.Event{Topic: "coord/update", Content: message.Coords})
			view := s.removeSelf(message.View)
			if message.Type == Reply {
				s.sampling.OnReceiveReply(NewView(message.Capacity, view))
			} else {
				str := addr.String()
				i := strings.LastIndex(str, ":")
				source := shared.Node{
					Ip:   str[:i],
					Port: 5050,
				}
				s.sampling.OnReceiveRequest(NewView(message.Capacity, view), source)
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
	case Reply:
	case Request:
		return message, nil
	default:
		log.Printf("Unknown message type= %d\n", message.Type)
	}
	return PViewMessage{}, err
}
