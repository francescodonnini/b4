package gossip

import (
	"b4/shared"
	"sync"
	"time"
)

// RetentionPolicy elimina le coordinate che non sono state aggiornate da più di un tempo pari a retention.
type RetentionPolicy struct {
	mu          *sync.RWMutex
	store       Store
	retention   time.Duration
	lastUpdates map[shared.Node]time.Time
}

func NewRetentionPolicy(store Store, retention time.Duration) *RetentionPolicy {
	return &RetentionPolicy{
		mu:          &sync.RWMutex{},
		store:       store,
		retention:   retention,
		lastUpdates: make(map[shared.Node]time.Time),
	}
}

func (r *RetentionPolicy) Remove(node shared.Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store.Remove(node)
	delete(r.lastUpdates, node)
}

func (r *RetentionPolicy) Peers() []shared.Node {
	r.removeOldNodes()
	return r.store.Peers()
}

func (r *RetentionPolicy) Items() []RemoteCoord {
	r.removeOldNodes()
	return r.store.Items()
}

// removeOldNodes rimuove tutti i nodi che non sono stati più aggiorni per più di un tempo pari a retention.
func (r *RetentionPolicy) removeOldNodes() {
	expired := make([]shared.Node, 0)
	for _, peer := range r.store.Peers() {
		if r.checkExpiration(peer) {
			expired = append(expired, peer)
		}
	}
	for _, peer := range expired {
		r.Remove(peer)
	}
}

// checkExpiration controlla se un certo nodo è stato aggiornato gli ultimi retention (ms).
// Ritorna true se il nodo è scaduto, falso altrimenti.
func (r *RetentionPolicy) checkExpiration(node shared.Node) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lastUpdate, ok := r.lastUpdates[node]
	expired := false
	if ok {
		expired = time.Now().Sub(lastUpdate).Milliseconds() > r.retention.Milliseconds()
	}
	return expired
}

// updateExpiration aggiorna l'ultima volta che una coordinata di un certo nodo è stata scritta.
func (r *RetentionPolicy) updateExpiration(node shared.Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastUpdates[node] = time.Now()
}

func (r *RetentionPolicy) Read(node shared.Node) (RemoteCoord, bool) {
	if r.checkExpiration(node) {
		r.Remove(node)
		return RemoteCoord{}, false
	}
	return r.store.Read(node)
}

func (r *RetentionPolicy) Save(coord RemoteCoord) {
	r.store.Save(coord)
	r.updateExpiration(coord.Owner)
}
