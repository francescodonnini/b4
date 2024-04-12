package gossip

import (
	"b4/shared"
	"math/rand"
	"sync"
	"time"
)

type Protocol interface {
	OnTimeout()
	OnReceiveReply(reply PViewMessage)
	OnReceiveRequest(request PViewMessage)
}

type Client interface {
	SendRequest(request PViewMessage, dest shared.Node)
	SendReply(reply PViewMessage, dest shared.Node)
}

type Impl struct {
	mu     sync.RWMutex
	id     shared.Node
	view   *PView
	client Client
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

func (i *Impl) OnReceiveReply(reply PViewMessage) {
	i.updateView(NewView(reply.Capacity, reply.View))
}

func (i *Impl) OnReceiveRequest(request PViewMessage) {
	v := i.view.Add(NewDescriptor(i.id, 0))
	reply := NewReply(v.capacity, v.Descriptors(), nil, request.Timestamp, i.id)
	i.client.SendReply(reply, request.Srv)
	i.updateView(NewView(request.Capacity, request.View))
}

func (i *Impl) OnTimeout() {
	p := i.view.GetDescriptor()
	v := i.view.Add(NewDescriptor(i.id, 0))
	request := NewRequest(v.capacity, v.Descriptors(), nil, time.Now(), i.id)
	i.client.SendRequest(request, p.Node)
}

func (i *Impl) updateView(other *PView) {
	i.mu.Lock()
	defer i.mu.Unlock()
	v := other.Merge(i.view)
	v = v.Select()
	v = v.Increase()
	i.view = v
}
