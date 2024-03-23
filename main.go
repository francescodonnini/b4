package main

import (
	"b4/discovery"
	"b4/sampling"
	"b4/shared"
	"b4/vivaldi"
	"b4/vivaldi/vivaldi_grpc"
	"b4/vivaldi/vivaldi_grpc/vivaldi_pb"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	ip := shared.GetIp()
	id := shared.Node{
		Ip:   ip.String(),
		Port: 5050,
	}
	endpoint := shared.Node{
		Ip:   "10.0.0.253",
		Port: 5050,
	}
	disc := discovery.NewDiscoveryService(endpoint, id)
	peers := disc.GetNodes()
	for range time.Tick(3 * time.Second) {
		if len(peers) >= 10 {
			break
		}
		peers = disc.GetNodes()
	}

	rand.NewSource(time.Now().Unix())
	rand.Shuffle(len(peers), func(i, j int) {
		peers[i], peers[j] = peers[j], peers[i]
	})
	s := grpc.NewServer()
	model := vivaldi.DefaultModel()
	registerServices(s, model)

	// lis, err := net.Listen("tcp", id.Address())
	// if err != nil {
	// 	log.Fatalf("Failed to listen: %s\n", err)
	// }
	// go startServer(s, lis)
	protocol := sampling.NewProtocol(id, 4, peers)
	go startUdpServer(context.Background(), id.Ip, id.Port, protocol)
	go startUdpClient(protocol)
	time.Sleep(600 * time.Second)
}

func startVivaldiService(client vivaldi.Client) {
	for range time.Tick(3 * time.Second) {
		client.Update()
	}
}

func registerServices(s *grpc.Server, model vivaldi.Model) {
	vivaldi_pb.RegisterVivaldiServer(s, vivaldi_grpc.NewServer(model))
}

func startServer(s *grpc.Server, lis net.Listener) {
	err := s.Serve(lis)
	if err != nil {
		log.Fatalf("Cannot serve: %s\n", err)
	}
}

func startUdpServer(ctx context.Context, ip string, port int, protocol sampling.Protocol) {
	address := fmt.Sprintf("%s:%d", ip, port)
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
			if message.Type == sampling.Reply {
				protocol.OnReceiveReply(sampling.NewView(message.Capacity, message.View))
			} else {
				s := addr.String()
				i := strings.LastIndex(s, ":")
				port, _ = strconv.Atoi(s[i+1:])
				source := shared.Node{
					Ip:   s[:i],
					Port: port,
				}
				protocol.OnReceiveRequest(sampling.NewView(message.Capacity, message.View), source)
			}
		}
	}()
}

func decodeMessage(payload []byte) (sampling.PViewMessage, error) {
	var message sampling.PViewMessage
	dec := gob.NewDecoder(bytes.NewReader(payload))
	err := dec.Decode(&message)
	if err != nil {
		log.Printf("Cannot decode request. Error: %s\n", err)
		return sampling.PViewMessage{}, err
	}
	switch message.Type {
	case sampling.Reply:
	case sampling.Request:
		return message, nil
	default:
		log.Printf("unknown message type= %d\n", message.Type)
	}
	return sampling.PViewMessage{}, err
}

func startUdpClient(protocol sampling.Protocol) {
	ticker := time.NewTicker(3 * time.Second)
	for range ticker.C {
		protocol.OnTimeout()
	}
}
