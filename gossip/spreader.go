package gossip

import (
	"b4/shared"
	event_bus "github.com/francescodonnini/pubsub"
	"sync"
)

type Value struct {
	RemoteCoord
	Counter int
}

func NewValue(coord RemoteCoord, counter int) *Value {
	return &Value{coord, counter}
}

func (v *Value) Dec() {
	v.Counter--
}

// Spreader gestisce il gossiping delle coordinate. La strategia utilizzata è detta feedback/counter.
// A ogni coordinata ricevuta (tramite il metodo Spread()) viene associato un contatore di valore maxN. Il contatore
// viene decrementato ogni volta che si riceve la stessa coordinata (sempre tramite Spread()), quando raggiunge valore
// zero la coordinata viene dimenticata. Una coordinata con contatore positivo ha possibilità di essere selezionata per
// il gossiping tramite il metodo Select().
type Spreader struct {
	mu *sync.RWMutex
	// Go non fornisce operazioni built-in per iterare sulle chiavi di map, inoltre l'ordine delle coppie chiave-valore
	// è implementation-dependent quindi si utilizza un array ausiliario per tenere le chiavi delle entry in cache.
	keySet    []shared.Node
	lastRound int
	cache     map[shared.Node]*Value
	bus       *event_bus.EventBus
	maxN      int
}

func NewSpreader(bus *event_bus.EventBus, maxN int) *Spreader {
	return &Spreader{
		bus:       bus,
		keySet:    make([]shared.Node, 0),
		cache:     make(map[shared.Node]*Value),
		mu:        &sync.RWMutex{},
		maxN:      maxN,
		lastRound: 0,
	}
}

// Select restituisce una lista di coordinate da diffondere (se ce ne sono). Il client deve specificare
// il numero massimo di coordinate che vuole diffondere, se il numero specificato è maggiore di quello delle coordinate
// che si conoscono allora vengono tornate tutte, altrimenti se ne tornano n selezionate a caso.
func (s *Spreader) Select(n int) []RemoteCoord {
	s.mu.Lock()
	defer s.mu.Unlock()
	var keys []shared.Node
	if n < len(s.cache) {
		keys = s.keySet[s.lastRound:]
	} else {
		keys = s.keySet[:]
		s.lastRound = 0
	}
	values := make([]RemoteCoord, 0)
	for _, k := range keys {
		values = append(values, s.cache[k].RemoteCoord)
		n--
	}
	s.lastRound += n
	if s.lastRound > len(s.cache) {
		s.lastRound = 0
	}
	return values
}

// Spread serve ad aggiungere le coordinate allo spreader in modo tale che queste possano essere diffuse e che
// si possa aggiornare il contatore.
func (s *Spreader) Spread(updates ...RemoteCoord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, coord := range updates {
		c, ok := s.cache[coord.Owner]
		if !ok || c.Age < coord.Age {
			if ok {
				s.keySet = shared.RemoveIf(s.keySet, func(node shared.Node) bool {
					return node == c.Owner
				})
			}
			s.keySet = append(s.keySet, coord.Owner)
			s.cache[coord.Owner] = NewValue(coord, s.maxN)
			s.bus.Publish(event_bus.Event{
				Topic:   "coord/store",
				Content: coord,
			})
		} else {
			c.Dec()
			if c.Counter == 0 {
				s.keySet = shared.RemoveIf(s.keySet, func(node shared.Node) bool {
					return node == c.Owner
				})
				delete(s.cache, c.Owner)
			}
		}
	}
}
