package gossip

import (
	logging "b4/logging"
	"b4/shared"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Protocol interface {
	OnTimeout()
	OnReceiveReply(reply Message)
	OnReceiveRequest(request Message)
}

type Client interface {
	SendRequest(view PViewMessage, dest shared.Node)
	SendReply(view PViewMessage, timestamp time.Time, dest shared.Node)
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

func (i *Impl) OnReceiveReply(reply Message) {
	i.updateView(NewView(reply.Capacity, reply.View))
}

func (i *Impl) OnReceiveRequest(request Message) {
	v := i.view.Add(NewDescriptor(i.id, 0))
	reply := PViewMessage{
		Capacity: v.capacity,
		View:     v.Descriptors(),
	}
	i.updateView(NewView(request.Capacity, request.View))
	i.client.SendReply(reply, request.Timestamp, request.Source)
}

func (i *Impl) OnTimeout() {
	p := i.view.GetDescriptor()
	v := i.view.Add(NewDescriptor(i.id, 0))
	request := PViewMessage{
		Capacity: v.capacity,
		View:     v.Descriptors(),
	}
	i.client.SendRequest(request, p.Node)
}

func (i *Impl) updateView(other *PView) {
	i.mu.Lock()
	defer i.mu.Unlock()
	v := other.Merge(i.view)
	v = v.Select()
	v = v.Increase()
	i.view = v
	logging.GetInstance("view").Println(str(v))
}

func str(v *PView) string {
	view := make([]string, 0)
	for _, n := range v.Descriptors() {
		view = append(view, n.Ip)
	}
	return strings.Join(view, ",")
}
