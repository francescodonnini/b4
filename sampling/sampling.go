package shared

import (
	"fmt"
	"math/rand"
	"time"
)

type PeerSamplingService interface {
	GetRandom() (string, bool)
}

type MockSampling struct {
	peers []string
}

func NewSamplingService() PeerSamplingService {
	peers := make([]string, 0)
	for i := 0; i < 16; i++ {
		peers = append(peers, fmt.Sprintf("10.0.0.%d", i))
	}
	return &MockSampling{peers: peers}
}

func (m MockSampling) GetRandom() (string, bool) {
	rand.NewSource(time.Now().Unix())
	n := len(m.peers)
	if n == 0 {
		return "", false
	}
	return m.peers[rand.Intn(n)], true
}
