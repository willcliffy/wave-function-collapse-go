package models

import uuid "github.com/satori/go.uuid"

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
		rotated = ModuleSlots{ms.NegativeZ, ms.PositiveY, ms.PositiveX, ms.PositiveZ, ms.NegativeY, ms.NegativeX}
	case 180:
		rotated = ModuleSlots{ms.NegativeX, ms.PositiveY, ms.NegativeZ, ms.PositiveX, ms.NegativeY, ms.PositiveZ}
	case 270:
		rotated = ModuleSlots{ms.PositiveZ, ms.PositiveY, ms.NegativeX, ms.NegativeZ, ms.NegativeY, ms.PositiveX}
	default: // 0 and multiples of 360
		rotated = ms
	}

	return rotated
}

type Module struct {
	ID       string
	Filename string
	Position Vector3
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

// Creates a new set of modules from a single file
// currently, this only supports rotations along the Y axis
func NewModuleSet(filename string, slots ModuleSlots) []Module {
	var modules []Module
	for _, rotation := range AllYRotations {
		modules = append(modules, NewModule(filename, slots.RotatedSlots(rotation.Y), rotation))
	}
	return modules
}
