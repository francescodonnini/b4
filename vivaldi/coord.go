package vivaldi

type Coord struct {
	Point Vec3d
}

func (c Coord) Add(other Coord) Coord {
	return Coord{
		Point: c.Point.Add(other.Point),
	}
}

func (c Coord) Magnitude() float64 {
	return c.Point.Magnitude()
}

func (c Coord) Scale(s float64) Coord {
	return Coord{
		Point: c.Point.Scale(s),
	}
}

func (c Coord) Sub(other Coord) Coord {
	return Coord{
		Point: c.Point.Sub(other.Point),
	}
}

func (c Coord) Unit() Coord {
	return Coord{
		Point: c.Point.Unit(),
	}
}
