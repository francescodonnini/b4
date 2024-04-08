package gossip

import (
	"b4/shared"
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

type UdpClient struct {
	spreader *Spreader
}

func NewClient(spreader *Spreader) Client {
	return &UdpClient{spreader: spreader}
}

// Send invia la vista parziale che il nodo ha sul sistema al peer dest. Oltre alla vista parziale vengono inviate
// le coordinate del mittente e/o quelle apprese da altri nodi (vedere spreader.go per informazioni su come vengono
// selezionate le coordinate da diffondere).
func (c *UdpClient) Send(request *PView, dest shared.Node) {
	conn, err := net.Dial("udp", dest.Address())
	if err != nil {
		log.Printf("Cannot dial to %s. Error: %s\n", dest.Address(), err)
		return
	}
	coords := c.spreader.Select(16)
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
