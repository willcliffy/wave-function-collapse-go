package models

type WFCMapMeta struct{}

type WFCMapData struct{}

type WFCMapCell struct {
	ModuleID string
	Filename string
	Position Vector3i
	Rotation int
}

type WFCMap []WFCMapCell
