package vivaldi

import "math"

type Coord struct {
	Point  Vec3d
	Height float64
}

func (c Coord) Add(other Coord) Coord {
	return Coord{
		Point:  c.Point.Add(other.Point),
		Height: c.Height + other.Height,
	}
}

func (c Coord) Magnitude() float64 {
	return c.Point.Magnitude() + c.Height
}

func (c Coord) Scale(s float64) Coord {
	return Coord{
		Point:  c.Point.Scale(s),
		Height: c.Height * math.Abs(s),
	}
}

func (c Coord) Sub(other Coord) Coord {
	return Coord{
		Point:  c.Point.Sub(other.Point),
		Height: c.Height + other.Height,
	}
}

func (c Coord) Unit() Coord {
	return Coord{
		Point:  c.Point.Unit(),
		Height: c.Height,
	}
}
