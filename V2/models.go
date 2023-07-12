package main

import "math"

// Note that the following "constants" are based on their corresponding Godot definitions
// See https://docs.godotengine.org/en/4.1/classes/class_vector3.html#constants
var (
	V3i_ZERO = Vector3i{}
	V3i_ONE  = Vector3i{1, 1, 1}

	V3i_LEFT  = Vector3i{-1, 0, 0}
	V3i_RIGHT = Vector3i{1, 0, 0}

	V3i_UP   = Vector3i{0, 1, 0}
	V3i_DOWN = Vector3i{0, -1, 0}

	V3i_FORWARD = Vector3i{0, 0, -1}
	V3i_BACK    = Vector3i{0, 0, 1}
)

type Vector3i struct {
	X, Y, Z int
}

func (v Vector3i) Add(other Vector3i) Vector3i {
	return Vector3i{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
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

type WFCPrototype struct {
	MeshName        string     `json:"mesh_name"`
	MeshRotation    int        `json:"mesh_rotation"`
	PosX            string     `json:"posX"`
	NegX            string     `json:"negX"`
	PosY            string     `json:"posY"`
	NegY            string     `json:"negY"`
	PosZ            string     `json:"posZ"`
	NegZ            string     `json:"negZ"`
	ConstrainTo     string     `json:"constrain_to"`
	ConstrainFrom   string     `json:"constrain_from"`
	Weight          int        `json:"weight"`
	ValidNeighbours [][]string `json:"valid_neighbours"`
}

type WFCMap struct {
	Size       Vector3i
	Prototypes [][][]WFCPrototype
}
