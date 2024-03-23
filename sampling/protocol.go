package sampling

import (
	"b4/shared"
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"sync"
)

type Protocol interface {
	OnTimeout()
	OnReceiveReply(request *PView)
	OnReceiveRequest(reply *PView, source shared.Node)
}

type Impl struct {
	mu      sync.RWMutex
	id      shared.Node
	view    *PView
	counter int64
}

func NewProtocol(id shared.Node, capacity int, peers []shared.Node) *Impl {
	view := make([]Descriptor, 0)
	for _, node := range peers {
		view = append(view, Descriptor{node, 0})
	}
	return &Impl{id: id, view: NewView(capacity, view)}
}

func (i *Impl) OnReceiveReply(reply *PView) {
	i.updateView(reply)
}

func (i *Impl) OnReceiveRequest(request *PView, source shared.Node) {
	reply := i.view.Add(NewDescriptor(i.id, i.now()))
	i.send(reply, source)
	i.updateView(request)
}

func (i *Impl) OnTimeout() {
	p := i.view.GetDescriptor()
	request := i.view.Add(NewDescriptor(i.id, i.now()))
	i.send(request, p.Node)
}

func (i *Impl) updateView(other *PView) {
	i.mu.Lock()
	defer i.mu.Unlock()
	v := other.Merge(i.view)
	i.view = v.Select()
}

func (i *Impl) send(request *PView, dest shared.Node) {
	conn, err := net.Dial("udp", dest.Address())
	if err != nil {
		log.Printf("Cannot dial to %s. Error: %s\n", dest.Address(), err)
		return
	}
	message := PViewMessage{
		Type:     Request,
		Capacity: request.Capacity(),
		View:     request.Descriptors(),
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

func (i *Impl) now() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.counter
}
