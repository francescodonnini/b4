package shared

import "sync"

type Publisher[K interface{}] struct {
	mu          sync.RWMutex
	subscribers map[string][]chan K
	closed      bool
}

func NewPublisher[K interface{}]() *Publisher[K] {
	return &Publisher[K]{
		subscribers: make(map[string][]chan K),
	}
}

func (p *Publisher[K]) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.closed {
		p.closed = true
		for _, sub := range p.subscribers {
			for _, ch := range sub {
				close(ch)
			}
		}
	}
}

func (p *Publisher[K]) Subscribe(topic string) <-chan K {
	p.mu.Lock()
	defer p.mu.Unlock()
	ch := make(chan K, 1)
	p.subscribers[topic] = append(p.subscribers[topic], ch)
	return ch
}

func (p *Publisher[K]) Publish(topic string, message K) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.closed {
		return
	}
	for _, ch := range p.subscribers[topic] {
		go func(ch chan K) {
			ch <- message
		}(ch)
	}
}
