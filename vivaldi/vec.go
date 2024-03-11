package vivaldi

import (
	"math"
)

type Vec3d [3]float64

func (v Vec3d) Add(other Vec3d) Vec3d {
	return Vec3d{v[0] + other[0], v[1] + other[1], v[2] + other[2]}
}

func (v Vec3d) Magnitude() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v Vec3d) Neg() Vec3d {
	return v.Scale(-1.0)
}

func (v Vec3d) Scale(s float64) Vec3d {
	return Vec3d{s * v[0], s * v[1], s * v[2]}
}

func (v Vec3d) Unit() Vec3d {
	m := v.Magnitude()
	return Vec3d{v[0] / m, v[1] / m, v[2] / m}
}
