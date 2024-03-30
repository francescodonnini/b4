package vivaldi

import (
	"math"
	"math/rand"
	"time"
)

type Vec3d [3]float64

func NewVec3d(x, y, z float64) Vec3d {
	return Vec3d{x, y, z}
}

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

func (v Vec3d) Sub(other Vec3d) Vec3d {
	return Vec3d{v[0] - other[0], v[1] - other[1], v[2] - other[2]}
}

func (v Vec3d) Unit() Vec3d {
	m := v.Magnitude()
	return Vec3d{v[0] / m, v[1] / m, v[2] / m}
}

func NewRandomUnit() Vec3d {
	rand.NewSource(time.Now().Unix())
	v := NewVec3d(rand.Float64(), rand.Float64(), rand.Float64())
	return v.Unit()
}
