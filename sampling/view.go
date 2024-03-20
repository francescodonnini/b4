package sampling

import (
	"math/rand"
	"sync"
	"time"
)

type PView interface {
	At(i int) Descriptor
	Capacity() int
	Descriptors() []Descriptor
	Increase() PView
	Length() int
	Merge(descriptors ...Descriptor) PView
	SelectPeer() (Descriptor, bool)
	SelectView() PView
}

type ViewImpl struct {
	c    int
	view []Descriptor
	mu   *sync.RWMutex
}

func NewView(c int, peers []Descriptor) PView {
	return &ViewImpl{
		c:    c,
		view: peers,
		mu:   &sync.RWMutex{},
	}
}

func (v ViewImpl) At(i int) Descriptor {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.view[i]
}

func (v ViewImpl) Capacity() int {
	return v.c
}

func (v ViewImpl) Descriptors() []Descriptor {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.view[:]
}

func (v ViewImpl) Increase() PView {
	view := make([]Descriptor, len(v.view))
	copy(view, v.Descriptors())
	for i := range view {
		view[i].Age++
	}
	return NewView(v.c, view)
}

func (v ViewImpl) Length() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return len(v.view)
}

func (v ViewImpl) Merge(descriptors ...Descriptor) PView {
	view := v.Descriptors()
	hits := make(map[string]Descriptor)
	for _, desc := range view {
		hits[desc.Address()] = desc
	}
	for _, desc := range descriptors {
		if d, ok := hits[desc.Address()]; ok {
			if desc.Age < d.Age {
				hits[d.Address()] = desc
			}
		}
	}
	view = make([]Descriptor, 0)
	for _, desc := range hits {
		view = append(view, desc)
	}
	return NewView(v.c, view)
}

func (v ViewImpl) SelectPeer() (Descriptor, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if len(v.view) == 0 {
		return Descriptor{}, false
	}
	rand.NewSource(time.Now().Unix())
	return v.view[rand.Intn(len(v.view))], true
}

func (v ViewImpl) SelectView() PView {
	if v.Length() == 0 {
		return NewView(0, make([]Descriptor, 0))
	}
	view := make([]Descriptor, len(v.view))
	copy(view, v.Descriptors())
	rand.NewSource(time.Now().Unix())
	rand.Shuffle(len(view), func(i, j int) {
		view[i], view[j] = view[j], view[i]
	})
	return NewView(v.c, view[:v.c])
}
