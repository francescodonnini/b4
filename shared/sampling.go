package shared

type PeerSampling interface {
	// GetRandom seleziona un peer a caso tra quelli che si conoscono del sistema. Se non si conosce alcun nodo
	// viene tornato false come secondo valore, true altrimenti.
	GetRandom() (Node, bool)
}
