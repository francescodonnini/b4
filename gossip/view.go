package gossip

import (
	"b4/shared"
	"math/rand"
	"sort"
	"sync"
)

// PView rappresenta le viste parziali che hanno i nodi sui membri del sistema. La vista ha una capacità (numero massimo
// di nodi che si conoscono) che è la stessa per tutti i nodi.
type PView struct {
	capacity int
	view     []Descriptor
	mu       *sync.RWMutex
}

func NewView(capacity int, view []Descriptor) *PView {
	return &PView{capacity: capacity, view: view, mu: &sync.RWMutex{}}
}

// Add si comporta essenzialmente come Merge ma aggiunge un solo descrittore, è un metodo di comodo che viene
// utilizzato dal nodo per aggiungere se stesso (con Age più recente).
func (v *PView) Add(descriptor Descriptor) *PView {
	view := v.Descriptors()
	i, exists := indexOf(view, descriptor.Node)
	if exists && view[i].Age < descriptor.Age {
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
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.view[rand.Intn(len(v.view))]
}

func (v *PView) Increase() *PView {
	v.mu.RLock()
	defer v.mu.RUnlock()
	view := make([]Descriptor, len(v.view))
	for i, desc := range v.view {
		view[i] = NewDescriptor(desc.Node, desc.Age+1)
	}
	return NewView(v.capacity, view)
}

// Merge ritorna una nuova lista che è il risultato dell'unione di v e view. Questo è l'unico caso in cui si produce
// una vista più grande della capacità, ci si aspetta che dopo Merge venga utilizzato Select per selezionare opportunamente
// i descrittori dall'unione delle due viste. Il merge non produce viste con descrittori duplicati (con stesso indirizzo). Nel caso di
// duplicati si tiene il descrittore con Age più basso.
func (v *PView) Merge(view *PView) *PView {
	set := make(map[string]Descriptor)
	for _, desc := range view.Descriptors() {
		set[desc.Address()] = desc
	}
	for _, desc := range v.Descriptors() {
		hit, ok := set[desc.Address()]
		if ok {
			if hit.Age > desc.Age {
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

// Select seleziona i primi c (capacità) nodi ordinati per Age.
func (v *PView) Select() *PView {
	view := v.Descriptors()
	rand.Shuffle(len(view), func(i, j int) {
		view[i], view[j] = view[j], view[i]
	})
	sort.Slice(view, func(i, j int) bool {
		return view[i].Age < view[j].Age
	})
	return NewView(v.capacity, view[:v.capacity])
}
