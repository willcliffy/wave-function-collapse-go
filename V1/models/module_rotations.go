package models

type ModuleRotation Vector3i

var (
	RotationX_000 = ModuleRotation{0, 0, 0}
	RotationX_090 = ModuleRotation{90, 0, 0}
	RotationX_180 = ModuleRotation{180, 0, 0}
	RotationX_270 = ModuleRotation{270, 0, 0}

	RotationY_000 = ModuleRotation{0, 0, 0}
	RotationY_090 = ModuleRotation{0, 90, 0}
	RotationY_180 = ModuleRotation{0, 180, 0}
	RotationY_270 = ModuleRotation{0, 270, 0}

	RotationZ_000 = ModuleRotation{0, 0, 0}
	RotationZ_090 = ModuleRotation{0, 0, 90}
	RotationZ_180 = ModuleRotation{0, 0, 180}
	RotationZ_270 = ModuleRotation{0, 0, 270}

	AllXRotations = []ModuleRotation{
		RotationX_000,
		RotationX_090,
		RotationX_180,
		RotationX_270,
	}

	AllYRotations = []ModuleRotation{
		RotationY_000,
		RotationY_090,
		RotationY_180,
		RotationY_270,
	}

	AllZRotations = []ModuleRotation{
		RotationZ_000,
		RotationZ_090,
		RotationZ_180,
		RotationZ_270,
	}
)
