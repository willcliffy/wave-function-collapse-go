package models_test

import (
	"fmt"
	"testing"

	"github.com/willcliffy/wfc_golang/models"
)

var (
	EmptySlot            = ""
	SquareSlot           = "square"
	BottomRectangleSlot  = "rect"
	PositiveRampUpSlot   = "rampUp"
	PositiveRampDownSlot = "rampDown"
	NegativeRampUpSlot   = PositiveRampDownSlot
	NegativeRampDownSlot = PositiveRampUpSlot

	RampSlots = models.ModuleSlots{
		PositiveX: SquareSlot,
		PositiveY: EmptySlot,
		PositiveZ: PositiveRampUpSlot,
		NegativeX: BottomRectangleSlot,
		NegativeY: SquareSlot,
		NegativeZ: NegativeRampDownSlot,
	}

	CubeSlots = models.ModuleSlots{
		PositiveX: SquareSlot,
		PositiveY: SquareSlot,
		PositiveZ: SquareSlot,
		NegativeX: SquareSlot,
		NegativeY: SquareSlot,
		NegativeZ: SquareSlot,
	}

	CornerSlots = models.ModuleSlots{
		PositiveX: PositiveRampUpSlot,
		PositiveY: EmptySlot,
		PositiveZ: BottomRectangleSlot,
		NegativeX: BottomRectangleSlot,
		NegativeY: SquareSlot,
		NegativeZ: NegativeRampDownSlot,
	}
)

func TestRotatedSlots(t *testing.T) {
	slots := models.ModuleSlots{
		PositiveX: "PositiveX",
		PositiveY: "PositiveY",
		PositiveZ: "PositiveZ",
		NegativeX: "NegativeX",
		NegativeY: "NegativeY",
		NegativeZ: "NegativeZ",
	}

	tests := []struct {
		name      string
		yRotation int
		want      models.ModuleSlots
	}{
		{
			name:      "Rotation 0",
			yRotation: 0,
			want:      slots,
		},
		{
			name:      "Rotation 90",
			yRotation: 90,
			want: models.ModuleSlots{
				PositiveX: "PositiveZ",
				PositiveY: "PositiveY",
				PositiveZ: "NegativeX",
				NegativeX: "NegativeZ",
				NegativeY: "NegativeY",
				NegativeZ: "PositiveX",
			},
		},
		{
			name:      "Rotation 180",
			yRotation: 180,
			want: models.ModuleSlots{
				PositiveX: "NegativeX",
				PositiveY: "PositiveY",
				PositiveZ: "NegativeZ",
				NegativeX: "PositiveX",
				NegativeY: "NegativeY",
				NegativeZ: "PositiveZ",
			},
		},
		{
			name:      "Rotation 270",
			yRotation: 270,
			want: models.ModuleSlots{
				PositiveX: "NegativeZ",
				PositiveY: "PositiveY",
				PositiveZ: "PositiveX",
				NegativeX: "PositiveZ",
				NegativeY: "NegativeY",
				NegativeZ: "NegativeX",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slots.RotatedSlots(tt.yRotation); !isEqual(got, tt.want) {
				t.Errorf("RotatedSlots() = %v, want %v", got, tt.want)
			}
		})
	}
}

func isEqual(a, b models.ModuleSlots) bool {
	return a.PositiveX == b.PositiveX &&
		a.PositiveY == b.PositiveY &&
		a.PositiveZ == b.PositiveZ &&
		a.NegativeX == b.NegativeX &&
		a.NegativeY == b.NegativeY &&
		a.NegativeZ == b.NegativeZ
}

func TestModule_CompatibleWith(t *testing.T) {
	tests := []struct {
		name     string
		moduleA  models.Module
		moduleB  models.Module
		expected bool
	}{
		{
			name: "Should return false when modules are more than one unit apart",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 2, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{},
			},
			expected: false,
		},
		{
			name: "Should return true when modules are one unit apart on X and have compatible slots",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{PositiveX: "slot1"},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 1, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{NegativeX: "slot1"},
			},
			expected: true,
		},
		{
			name: "Should return false when modules are one unit apart on X but have incompatible slots",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{PositiveX: "slot1"},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 1, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{NegativeX: "slot2"},
			},
			expected: false,
		},
		{
			name: "Should return false when modules are more than one unit apart",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 0, Y: 2, Z: 0},
				Slots:    models.ModuleSlots{},
			},
			expected: false,
		},
		{
			name: "Should return true when modules are one unit apart on Y and have compatible slots",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{PositiveY: "slot1"},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 0, Y: 1, Z: 0},
				Slots:    models.ModuleSlots{NegativeY: "slot1"},
			},
			expected: true,
		},
		{
			name: "Should return false when modules are one unit apart on Y but have incompatible slots",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{PositiveY: "slot1"},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 0, Y: 1, Z: 0},
				Slots:    models.ModuleSlots{NegativeY: "slot2"},
			},
			expected: false,
		},
		{
			name: "Should return false when modules are more than one unit apart",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 0, Y: 0, Z: 2},
				Slots:    models.ModuleSlots{},
			},
			expected: false,
		},
		{
			name: "Should return true when modules are one unit apart on Z and have compatible slots",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{PositiveZ: "slot1"},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 0, Y: 0, Z: 1},
				Slots:    models.ModuleSlots{NegativeZ: "slot1"},
			},
			expected: true,
		},
		{
			name: "Should return false when modules are one unit apart on Z but have incompatible slots",
			moduleA: models.Module{
				ID:       "1",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    models.ModuleSlots{PositiveZ: "slot1"},
			},
			moduleB: models.Module{
				ID:       "2",
				Position: models.Vector3i{X: 0, Y: 0, Z: 1},
				Slots:    models.ModuleSlots{NegativeZ: "slot2"},
			},
			expected: false,
		},
		{
			name: "Ramp and cube should be compatible when oriented correctly",
			moduleA: models.Module{
				ID:       "cube.glb",
				Position: models.Vector3i{X: 1, Y: 0, Z: 0},
				Slots:    CubeSlots,
				Rotation: models.RotationY_000,
			},
			moduleB: models.Module{
				ID:       "ramp.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    RampSlots,
				Rotation: models.RotationY_000,
			},
			expected: true,
		},
		{
			name: "Ramp and cube should be compatible when oriented correctly",
			moduleA: models.Module{
				ID:       "cube.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    CubeSlots,
				Rotation: models.RotationY_000,
			},
			moduleB: models.Module{
				ID:       "ramp.glb",
				Position: models.Vector3i{X: 1, Y: 0, Z: 0},
				Slots:    RampSlots.RotatedSlots(models.RotationY_180.Y),
				Rotation: models.RotationY_180,
			},
			expected: true,
		},
		{
			name: "Ramp and cube should be incompatible when oriented incorrectly",
			moduleA: models.Module{
				ID:       "cube.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    CubeSlots,
				Rotation: models.RotationY_000,
			},
			moduleB: models.Module{
				ID:       "ramp.glb",
				Position: models.Vector3i{X: 1, Y: 0, Z: 0},
				Slots:    RampSlots,
				Rotation: models.RotationY_000,
			},
			expected: false,
		},
		{
			name: "Ramp and cube should be incompatible when oriented incorrectly",
			moduleA: models.Module{
				ID:       "cube.glb",
				Position: models.Vector3i{X: 1, Y: 0, Z: 0},
				Slots:    CubeSlots,
				Rotation: models.RotationY_000,
			},
			moduleB: models.Module{
				ID:       "ramp.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    RampSlots.RotatedSlots(models.RotationY_090.Y),
				Rotation: models.RotationY_090,
			},
			expected: false,
		},
		{
			name: "Ramp and cube should not be compatible when oriented incorrectly",
			moduleA: models.Module{
				ID:       "cube.glb",
				Position: models.Vector3i{X: 1, Y: 0, Z: 0},
				Slots:    CubeSlots,
				Rotation: models.RotationY_000,
			},
			moduleB: models.Module{
				ID:       "ramp.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    RampSlots.RotatedSlots(models.RotationY_270.Y),
				Rotation: models.RotationY_270,
			},
			expected: false,
		},
		{
			name: "Ramp and corner should be compatible when oriented correctly",
			moduleA: models.Module{
				ID:       "corner.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 1},
				Slots:    CornerSlots,
				Rotation: models.RotationY_000,
			},
			moduleB: models.Module{
				ID:       "ramp.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    RampSlots,
				Rotation: models.RotationY_000,
			},
			expected: true,
		},
		{
			name: "Ramp and corner should be compatible when oriented correctly",
			moduleA: models.Module{
				ID:       "corner.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 1},
				Slots:    CornerSlots.RotatedSlots(models.RotationY_090.Y),
				Rotation: models.RotationY_090,
			},
			moduleB: models.Module{
				ID:       "ramp.glb",
				Position: models.Vector3i{X: 0, Y: 0, Z: 0},
				Slots:    RampSlots.RotatedSlots(models.RotationY_180.Y),
				Rotation: models.RotationY_180,
			},
			expected: true,
		},
		{
			name: "Ramp and corner should be incompatible when oriented incorrectly",
			moduleA: models.Module{
				ID:       "corner.glb",
				Position: models.Vector3i{X: 6, Y: 0, Z: 5},
				Slots:    CornerSlots.RotatedSlots(models.RotationY_000.Y),
				Rotation: models.RotationY_000,
			},
			moduleB: models.Module{
				ID:       "ramp.glb",
				Position: models.Vector3i{X: 7, Y: 0, Z: 5},
				Slots:    RampSlots.RotatedSlots(models.RotationY_270.Y),
				Rotation: models.RotationY_270,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.moduleA.CompatibleWith(&tt.moduleB); got != tt.expected {
				t.Errorf("Module.CompatibleWith() = %v, want %v", got, tt.expected)
				fmt.Printf("%v\n  ModA: %+v\n  ModB: %+v\n", tt.name, tt.moduleA, tt.moduleB)
			}
		})
	}
}
