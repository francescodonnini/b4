package vivaldi

import (
	"math"
	"math/rand"
)

type Vec []float64

func (v Vec) Add(other Vec) Vec {
	checkDim(v, other)
	values := make([]float64, v.Dim())
	for i := 0; i < v.Dim(); i++ {
		values[i] = v[i] + other[i]
	}
	return values
}

func checkDim(u, v Vec) {
	if u.Dim() != v.Dim() {
		panic("dimensions do not match!")
	}
}

func (v Vec) Dim() int {
	return len(v)
}

func (v Vec) Magnitude() float64 {
	sum := 0.0
	for _, x := range v {
		sum += x * x
	}
	return math.Sqrt(sum)
}

func (v Vec) Neg() Vec {
	return v.Scale(-1.0)
}

func (v Vec) Scale(s float64) Vec {
	values := make([]float64, v.Dim())
	for i, x := range v {
		values[i] = x * s
	}
	return values
}

func (v Vec) Sub(other Vec) Vec {
	return v.Add(other.Scale(-1.0))
}

func (v Vec) Unit() Vec {
	m := v.Magnitude()
	return v.Scale(1 / m)
}

func NewRandomUnit(n int) Vec {
	values := make([]float64, n)
	for i := 0; i < n; i++ {
		values[i] = rand.Float64()
	}
	return values
}
