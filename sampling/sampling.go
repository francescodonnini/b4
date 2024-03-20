package sampling

import (
	"b4/shared"
	"log"
	"math/rand"
	"time"
)

type PeerSampling interface {
	GetRandom() (shared.Node, bool)
}

type PeerSamplingProtocol interface {
	OnTimeout()
	OnReceive(view PView) PView
}

type PeerSamplingService struct {
	id     shared.Node
	View   PView
	client Client
}

func NewSamplingService(id shared.Node, c int, peers []shared.Node, client Client) PeerSamplingService {
	view := make([]Descriptor, 0)
	for _, p := range peers {
		view = append(view, FreshDescriptor(p))
	}
	return PeerSamplingService{
		id:     id,
		View:   NewView(c, view),
		client: client,
	}
}

func (p PeerSamplingService) GetRandom() (shared.Node, bool) {
	if p.View.Length() == 0 {
		return shared.Node{}, false
	}
	rand.NewSource(time.Now().Unix())
	desc := p.View.At(rand.Intn(p.View.Length()))
	return desc.Node, true
}

func (p PeerSamplingService) OnTimeout() {
	peer, ok := p.View.SelectPeer()
	if !ok {
		log.Printf("On timeout: no peer selected!")
		return
	}
	buffer := p.View.Merge(Descriptor{
		Node: p.id,
		Age:  0,
	})
	view, err := p.client.Exchange(buffer, peer.Node)
	if err != nil {
		return
	}
	view = view.Increase()
	buffer = view.Merge(buffer.Descriptors()...)
	p.View = buffer.SelectView()
}

func (p PeerSamplingService) OnReceive(view PView) PView {
	view = view.Increase()
	buffer := p.View.Merge(FreshDescriptor(p.id))
	v := p.View.Merge(view.Descriptors()...)
	p.View = v.SelectView()
	return buffer
}
