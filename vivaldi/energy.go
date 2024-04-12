package vivaldi

import (
	"b4/shared"
	eventbus "github.com/francescodonnini/pubsub"
	"time"
)

// EnergySlidingWindow
// Le coordinate di un nodo cambiano continuamente, molte applicazioni non hanno bisogno di ricevere ogni valore della
// coordinata calcolato. Si distingue infatti tra coordinate di applicazione e coordinate di sistema. Le prime tendono
// appunto a rimanere più stabili. Uno dei metodi per limitare il rate di aggiornamento delle coordinate è utilizzare
// un protocollo a finestra scorrevole (vedere http://nrs.harvard.edu/urn-3:HUL.InstRepos:25686820).
// startWindow e currentWindow hanno la stessa dimensione. Inizialmente viene riempita solo la prima, quando la prima diventa piena,
// si copiano tutti i valori nella seconda. Da questo punto in poi, fino a quando energyDistance > threshold si elimina
// il primo valore di currentWindow e se ne aggiunge uno nuovo alla fine (da qui finestra scorrevole).
type EnergySlidingWindow struct {
	startWindow   []Coord
	currentWindow []Coord
	windowSize    int
	threshold     float64
	ca            Coord
	bus           *eventbus.EventBus
}

func NewEnergySlidingWindow(windowSize int, threshold float64, bus *eventbus.EventBus) *EnergySlidingWindow {
	return &EnergySlidingWindow{windowSize: windowSize, threshold: threshold, startWindow: make([]Coord, 0), currentWindow: make([]Coord, 0), bus: bus}
}

func (e *EnergySlidingWindow) Update(coord Coord) {
	if len(e.startWindow) < e.windowSize {
		e.startWindow = append(e.startWindow, coord)
		if len(e.startWindow) == e.windowSize {
			e.currentWindow = make([]Coord, e.windowSize)
			copy(e.currentWindow, e.startWindow)
		}
	} else {
		e.currentWindow = append(e.currentWindow[1:], coord)
		if energyDistance(e.startWindow, e.currentWindow) > e.threshold {
			m := len(e.currentWindow) / 2
			e.ca = e.currentWindow[m]
			e.startWindow = e.startWindow[:0]
			e.currentWindow = e.currentWindow[:0]
			e.bus.Publish(eventbus.Event{
				Topic:   "coord/app",
				Content: shared.NewPair(e.ca, time.Now()),
			})
		}
	}
}

func energyDistance(set1, set2 []Coord) float64 {
	n1 := float64(len(set1))
	n2 := float64(len(set2))
	s1 := 2.0 / (n1 * n2) * sumOfDistances(set1, set2)
	s2 := 1.0 / (n1 * n1) * sumOfDistances(set1, set1)
	s3 := 1.0 / (n2 * n2) * sumOfDistances(set2, set2)
	return (n1 * n2) / (n1 + n2) * (s1 - s2 - s3)
}

func sumOfDistances(set1, set2 []Coord) float64 {
	s := 0.0
	for _, a := range set1 {
		for _, b := range set2 {
			s += a.Sub(b).Magnitude()
		}
	}
	return s
}
