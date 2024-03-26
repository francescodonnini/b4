package main

import (
	sampl "b4/gossip"
	"b4/shared"
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"strings"
)

type UdpServer struct {
	id       shared.Node
	sampling sampl.Protocol
	bus      *shared.EventBus
}

func NewUdpServer(id shared.Node, sampling sampl.Protocol, bus *shared.EventBus) *UdpServer {
	return &UdpServer{id: id, sampling: sampling, bus: bus}
}

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
		buf := make([]byte, 4096)
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
			view := s.removeSelf(message.View)
			s.bus.Publish(shared.Event{Topic: "coord/update", Content: message.Coords})
			if message.Type == sampl.Reply {
				s.sampling.OnReceiveReply(sampl.NewView(message.Capacity, view))
			} else {
				str := addr.String()
				i := strings.LastIndex(str, ":")
				port, _ := strconv.Atoi(str[i+1:])
				source := shared.Node{
					Ip:   str[:i],
					Port: port,
				}
				s.sampling.OnReceiveRequest(sampl.NewView(message.Capacity, view), source)
			}
		}
	}()
}

func (s *UdpServer) removeSelf(view []sampl.Descriptor) []sampl.Descriptor {
	return Filter(view, func(descriptor sampl.Descriptor) bool {
		return descriptor.Node == s.id
	})
}

func decodeMessage(payload []byte) (sampl.PViewMessage, error) {
	var message sampl.PViewMessage
	dec := gob.NewDecoder(bytes.NewReader(payload))
	err := dec.Decode(&message)
	if err != nil {
		log.Printf("Cannot decode request. Error: %s\n", err)
		return sampl.PViewMessage{}, err
	}
	switch message.Type {
	case sampl.Reply:
	case sampl.Request:
		return message, nil
	default:
		log.Printf("unknown message type= %d\n", message.Type)
	}
	return sampl.PViewMessage{}, err
}
