package vivaldi

type EnergySlidingWindow struct {
	startWindow   []Coord
	currentWindow []Coord
	windowSize    int
	threshold     float64
	ca            Coord
}

func NewEnergySlidingWindow(windowSize int, threshold float64) *EnergySlidingWindow {
	return &EnergySlidingWindow{windowSize: windowSize, threshold: threshold, startWindow: make([]Coord, 0), currentWindow: make([]Coord, 0)}
}

func (e *EnergySlidingWindow) Update(coord Coord) {
	if len(e.startWindow) < e.windowSize {
		e.startWindow = append(e.startWindow, coord)
		if len(e.startWindow) == e.windowSize {
			e.currentWindow = make([]Coord, e.windowSize)
			copy(e.currentWindow, e.startWindow)
		}
	}
	e.currentWindow = append(e.currentWindow[1:], coord)
	if energyDistance(e.startWindow, e.currentWindow) > e.threshold {
		m := len(e.currentWindow) / 2
		e.ca = e.currentWindow[m]
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
