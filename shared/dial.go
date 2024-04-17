package shared

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
)

type Dialer interface {
	Dial(target Node) (*grpc.ClientConn, error)
	Close()
}

type Impl struct {
	cache map[Node]*grpc.ClientConn
	mu    sync.RWMutex
}

func NewDialer() Dialer {
	return &Impl{cache: make(map[Node]*grpc.ClientConn)}
}

func (c *Impl) Dial(target Node) (*grpc.ClientConn, error) {
	c.mu.RLock()
	if conn, ok := c.cache[target]; ok {
		c.mu.RUnlock()
		return conn, nil
	}
	c.mu.RUnlock()
	conn, err := grpc.Dial(target.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	c.mu.Lock()
	defer c.mu.Unlock()
	if err != nil {
		log.Printf("Cannot connect to %s. Error: %s\n", target.Address(), err)
		return nil, err
	}
	c.cache[target] = conn
	return conn, nil
}

func (c *Impl) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, conn := range c.cache {
		err := conn.Close()
		if err != nil {
			log.Printf("cannot close grpc connection, error: %s\n", err)
			continue
		}
	}
}
