package sampling

import (
	"b4/shared"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type PView struct {
	capacity int
	view     []Descriptor
	mu       *sync.RWMutex
}

func NewView(capacity int, view []Descriptor) *PView {
	return &PView{capacity: capacity, view: view, mu: &sync.RWMutex{}}
}

func (v *PView) Add(descriptor Descriptor) *PView {
	view := v.Descriptors()
	i, exists := indexOf(view, descriptor.Node)
	if exists {
		view[i] = descriptor
	} else {
		view = append(view, descriptor)
	}
	return NewView(v.capacity, view)
}

func indexOf(descriptors []Descriptor, node shared.Node) (int, bool) {
	for i, it := range descriptors {
		if it.Node == node {
			return i, true
		}
	}
	return -1, false
}

func (v *PView) Capacity() int {
	return v.capacity
}

func (v *PView) Descriptors() []Descriptor {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.view[:]
}

func (v *PView) GetDescriptor() Descriptor {
	rand.NewSource(time.Now().Unix())
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.view[rand.Intn(len(v.view))]
}

func (v *PView) Merge(view *PView) *PView {
	set := make(map[string]Descriptor)
	for _, desc := range view.Descriptors() {
		set[desc.Address()] = desc
	}
	for _, desc := range v.Descriptors() {
		hit, ok := set[desc.Address()]
		if ok {
			if hit.Timestamp < desc.Timestamp {
				set[desc.Address()] = desc
			}
		} else {
			set[desc.Address()] = desc
		}
	}
	buffer := make([]Descriptor, 0)
	for _, desc := range set {
		buffer = append(buffer, desc)
	}

	return NewView(v.capacity, buffer)
}

func (v *PView) Select() *PView {
	view := v.Descriptors()
	sort.Slice(view, func(i, j int) bool {
		return view[i].Timestamp > view[j].Timestamp
	})
	return NewView(v.capacity, view[:v.capacity])
}
