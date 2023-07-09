package models_test

import (
	"testing"

	"github.com/willcliffy/wfc_golang/models"
)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.moduleA.CompatibleWith(&tt.moduleB); got != tt.expected {
				t.Errorf("Module.CompatibleWith() = %v, want %v", got, tt.expected)
			}
		})
	}
}
