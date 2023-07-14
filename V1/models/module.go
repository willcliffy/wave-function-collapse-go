package models

import (
	uuid "github.com/satori/go.uuid"
)

type ModuleSlots struct {
	PositiveX string
	PositiveY string
	PositiveZ string
	NegativeX string
	NegativeY string
	NegativeZ string
}

func (ms ModuleSlots) RotatedSlots(yRotation int) ModuleSlots {
	if yRotation%90 != 0 {
		panic("ModuleSlots.RotatedSlots only supports rotations 0, 90, 180, 270")
	}

	var rotated ModuleSlots

	yRotation = int(yRotation) % 360
	if yRotation < 0 {
		yRotation += 360
	}

	switch yRotation {
	case 90:
		rotated = ModuleSlots{ms.PositiveZ, ms.PositiveY, ms.NegativeX, ms.NegativeZ, ms.NegativeY, ms.PositiveX}
	case 180:
		rotated = ModuleSlots{ms.NegativeX, ms.PositiveY, ms.NegativeZ, ms.PositiveX, ms.NegativeY, ms.PositiveZ}
	case 270:
		rotated = ModuleSlots{ms.NegativeZ, ms.PositiveY, ms.PositiveX, ms.PositiveZ, ms.NegativeY, ms.NegativeX}
	default: // 0 and multiples of 360
		rotated = ms
	}

	return rotated
}

type Module struct {
	ID       string
	Filename string
	Position Vector3i
	Rotation ModuleRotation
	Slots    ModuleSlots
}

func NewModule(filename string, slots ModuleSlots, rotation ModuleRotation) Module {
	return Module{
		ID:       uuid.NewV4().String(),
		Filename: filename,
		Rotation: rotation,
		Slots:    slots,
	}
}

func (m Module) ToWFCMapCell() WFCMapCell {
	return WFCMapCell{
		ModuleID: m.ID,
		Filename: m.Filename,
		Position: m.Position,
		Rotation: m.Rotation.Y,
	}
}

// Creates a new set of modules from a single file
// currently, this only supports rotations along the Y axis
func NewModuleSet(filename string, slots ModuleSlots) []Module {
	var modules []Module
	for _, rotation := range AllYRotations {
		modules = append(modules, NewModule(filename, slots.RotatedSlots(rotation.Y), rotation))
	}
	return modules
}

func (m Module) Copy() *Module {
	return &Module{
		ID:       m.ID,
		Filename: m.Filename,
		Position: m.Position,
		Rotation: m.Rotation,
		Slots:    m.Slots,
	}
}

func (m Module) CompatibleWith(other *Module) bool {
	if other == nil {
		return false
	}

	if m.Position.Distance(other.Position) > 1 {
		return false
	}

	if m.Position.X != other.Position.X {
		if m.Position.X > other.Position.X {
			return m.Slots.NegativeX == other.Slots.PositiveX
		} else {
			return m.Slots.PositiveX == other.Slots.NegativeX
		}
	} else if m.Position.Y != other.Position.Y {
		if m.Position.Y > other.Position.Y {
			return m.Slots.NegativeY == other.Slots.PositiveY
		} else {
			return m.Slots.PositiveY == other.Slots.NegativeY
		}
	} else if m.Position.Z != other.Position.Z {
		if m.Position.Z > other.Position.Z {
			return m.Slots.NegativeZ == other.Slots.PositiveZ
		} else {
			return m.Slots.PositiveZ == other.Slots.NegativeZ
		}
	}

	return true
}
