package gossip

import (
	"b4/shared"
	"math/rand"
	"sync"
	"time"
)

type Protocol interface {
	OnTimeout()
	OnReceiveReply(request *PView)
	OnReceiveRequest(reply *PView, source shared.Node)
}

type Client interface {
	Send(view *PView, dest shared.Node)
}

type Impl struct {
	mu      sync.RWMutex
	id      shared.Node
	view    *PView
	counter int64
	client  Client
}

func (i *Impl) GetRandom() (shared.Node, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	peers := i.view.Descriptors()
	if len(peers) <= 0 {
		return shared.Node{}, false
	}
	rand.NewSource(time.Now().Unix())
	peer := peers[rand.Intn(len(peers))]
	return peer.Node, true
}

func NewProtocol(id shared.Node, capacity int, peers []shared.Node, client Client) *Impl {
	view := make([]Descriptor, 0)
	for _, node := range peers {
		view = append(view, Descriptor{node, 0})
	}
	return &Impl{id: id, view: NewView(capacity, view), client: client}
}

func (i *Impl) OnReceiveReply(reply *PView) {
	i.updateView(reply)
}

func (i *Impl) OnReceiveRequest(request *PView, source shared.Node) {
	reply := i.view.Add(NewDescriptor(i.id, i.now()))
	i.client.Send(reply, source)
	i.updateView(request)
}

func (i *Impl) OnTimeout() {
	p := i.view.GetDescriptor()
	request := i.view.Add(NewDescriptor(i.id, i.now()))
	i.client.Send(request, p.Node)
	i.incCounter()
}

func (i *Impl) updateView(other *PView) {
	i.mu.Lock()
	defer i.mu.Unlock()
	v := other.Merge(i.view)
	i.view = v.Select()
}

func (i *Impl) now() int64 {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.counter
}

func (i *Impl) incCounter() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.counter += 1
}
