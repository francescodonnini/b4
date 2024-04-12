package gossip

import (
	"b4/shared"
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"
)

type UdpClient struct {
	spreader *Spreader
	srv      shared.Node
}

func NewClient(spreader *Spreader, srv shared.Node) Client {
	return &UdpClient{spreader: spreader, srv: srv}
}

// SendRequest invia la vista parziale che il nodo ha sul sistema al peer dest. Oltre alla vista parziale vengono inviate
// le coordinate del mittente e/o quelle apprese da altri nodi (vedere spreader.go per informazioni su come vengono
// selezionate le coordinate da diffondere).
func (c *UdpClient) SendRequest(view PViewMessage, dest shared.Node) {
	conn, err := net.Dial("udp", dest.Address())
	if err != nil {
		log.Printf("Cannot dial to %s. Error: %s\n", dest.Address(), err)
		return
	}
	coords := c.spreader.Select(16)
	request := NewRequest(view, coords, time.Now(), c.srv)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(request)
	if err != nil {
		log.Printf("Cannot encode request. Error: %s\n", err)
		return
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Printf("Cannot send request to %s. Error: %s\n", dest.Address(), err)
	}
}

func (c *UdpClient) SendReply(view PViewMessage, timestamp time.Time, dest shared.Node) {
	conn, err := net.Dial("udp", dest.Address())
	if err != nil {
		log.Printf("Cannot dial to %s. Error: %s\n", dest.Address(), err)
		return
	}
	coords := c.spreader.Select(16)
	reply := NewReply(view, coords, timestamp, c.srv)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(reply)
	if err != nil {
		log.Printf("Cannot encode request. Error: %s\n", err)
		return
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Printf("Cannot send request to %s. Error: %s\n", dest.Address(), err)
	}
}
