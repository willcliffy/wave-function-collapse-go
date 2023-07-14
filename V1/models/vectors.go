package models

import "math"

type Vector3 struct {
	X, Y, Z float64
}

type Vector3i struct {
	X, Y, Z int
}

func (v Vector3i) Distance(other Vector3i) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z

	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

func (v Vector3i) Length() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z))
}

type Vector2 struct {
	X, Z float64
}
