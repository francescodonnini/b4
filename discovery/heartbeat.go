package discovery

import (
	"b4/shared"
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

type Beat struct{}

type HeartBeatClient struct {
	detector shared.Node
}

func NewHeartBeatClient(detector shared.Node) *HeartBeatClient {
	return &HeartBeatClient{detector: detector}
}

func (c *HeartBeatClient) Beat() {
	conn, err := net.Dial("udp", c.detector.Address())
	if err != nil {
		log.Printf("Cannot dial to %s. Error: %s\n", c.detector.Address(), err)
		return
	}
	beat := Beat{}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(beat)
	if err != nil {
		log.Printf("Cannot encode request. Error: %s\n", err)
		return
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Printf("Cannot send request to %s. Error: %s\n", c.detector.Address(), err)
	}
}
