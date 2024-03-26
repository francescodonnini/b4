package sampling

import (
	"b4/shared"
	"b4/vivaldi"
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"unsafe"
)

type UdpClient struct {
	spreader *Spreader
}

func NewClient(spreader *Spreader) Client {
	return &UdpClient{spreader: spreader}
}

func (c *UdpClient) Send(request *PView, dest shared.Node) {
	conn, err := net.Dial("udp", dest.Address())
	if err != nil {
		log.Printf("Cannot dial to %s. Error: %s\n", dest.Address(), err)
		return
	}
	descriptors := request.Descriptors()
	numOfBytes := unsafe.Sizeof(descriptors) + unsafe.Sizeof(request.Capacity()) + unsafe.Sizeof(Request)
	coords, err := c.spreader.Select(numOfBytes)
	if err != nil {
		coords = make(map[shared.Node]vivaldi.Coord)
	}
	message := PViewMessage{
		Type:     Request,
		Capacity: request.Capacity(),
		View:     request.Descriptors(),
		Coords:   coords,
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(message)
	if err != nil {
		log.Printf("Cannot encode request. Error: %s\n", err)
		return
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Printf("Cannot send request to %s. Error: %s\n", dest.Address(), err)
	}
}
