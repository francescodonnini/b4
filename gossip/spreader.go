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
	// L'ordine delle chiavi quando si itera su di esse tramite range è implementation-dependent quindi si utilizza un array
	// ausiliario per poter diffondere le coordinate secondo Round-Robin.
	keySet []shared.Node
	// lastRound indice dell'ultima coordinata selezionata (diffusa).
	lastRound int
	infected  map[shared.Node]*Value
	removed   map[shared.Node]RemoteCoord
	bus       *event_bus.EventBus
	// maxN è il numero massimo di volte che si è disposti a ricevere una coordinata che già si sta diffondendo prima
	// di smettere di diffondere del tutto la coordinata.
	maxN int
}

func NewSpreader(bus *event_bus.EventBus, maxN int) *Spreader {
	return &Spreader{
		bus:       bus,
		keySet:    make([]shared.Node, 0),
		infected:  make(map[shared.Node]*Value),
		removed:   make(map[shared.Node]RemoteCoord),
		mu:        &sync.RWMutex{},
		maxN:      maxN,
		lastRound: 0,
	}
}

// Select restituisce una lista di coordinate da diffondere (se ce ne sono). Il client deve specificare
// il numero massimo di coordinate che vuole diffondere, se il numero specificato è maggiore di quello delle coordinate
// che si conoscono allora vengono selezionate tutte, altrimenti se ne tornano n selezionate secondo una politica Round-Robin.
func (s *Spreader) Select(n int) []RemoteCoord {
	s.mu.Lock()
	defer s.mu.Unlock()
	var keys []shared.Node
	if n < len(s.infected) {
		keys = s.keySet[s.lastRound:]
	} else {
		keys = s.keySet[:]
		s.lastRound = 0
	}
	values := make([]RemoteCoord, 0)
	for _, k := range keys {
		values = append(values, s.infected[k].RemoteCoord)
	}
	s.lastRound += n
	if s.lastRound > len(s.infected) {
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
		// La coordinata non è attualmente nel gruppo delle coordinate che si vuole diffondere, però potrebbe
		// essere tra quelle che si ha diffuso in passato.
		if v, ok := s.infected[coord.Owner]; !ok {
			c, ok := s.removed[coord.Owner]
			// La coordinata non è stata mai diffusa in passato oppure
			// La coordinata è stata già diffusa in passato però quella appena ricevuto è una versione aggiornata
			if !ok || (ok && c.Age.Before(coord.Age)) {
				delete(s.removed, c.Owner)
				s.updateInfectedSet(coord)
			}
		} else if v.Age.Before(coord.Age) {
			// La coordinata sta venendo diffusa però si ha ricevuto una versione più recente di quella che si conosce allora
			// occorre aggiungerla alla lista delle coordinate da diffondere.
			// Si stava diffondendo una versione obsoleta della coordinata che occorre quindi eliminare dalla lista
			// delle coordinate da diffondere.
			if ok {
				s.keySet = removeByKey(s.keySet, v.Owner)
			}
			s.updateInfectedSet(coord)
		} else {
			// È stata ricevuta una coordinata che già si conosce (e si sta diffondendo) quindi si decrementa il contatore
			v.Dec()
			if v.Counter == 0 {
				s.keySet = removeByKey(s.keySet, v.Owner)
				delete(s.infected, v.Owner)
				s.removed[v.Owner] = v.RemoteCoord
			}
		}
	}
}

func (s *Spreader) updateInfectedSet(c RemoteCoord) {
	s.keySet = append(s.keySet, c.Owner)
	s.infected[c.Owner] = NewValue(c, s.maxN)
	s.bus.Publish(event_bus.Event{
		Topic:   "coord/store",
		Content: c,
	})
}

func removeByKey(keySet []shared.Node, key shared.Node) []shared.Node {
	return shared.RemoveIf(keySet, func(node shared.Node) bool {
		return node == key
	})
}
