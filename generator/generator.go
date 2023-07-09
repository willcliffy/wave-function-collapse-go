package generator

import "github.com/willcliffy/wfc_golang/models"

type IWFCMapGenerator interface {
	Initialize(allPossibleModules []models.Module) bool
	Run(done chan bool)
	observe_and_collapse() bool
	propagate() bool
}

type WFCMapGenerator struct {
	Size models.Vector3i

	// x, y, z -> list of potential modules
	CurrentMap [][][][]models.Module
}

func NewWFCMapGenerator(size models.Vector3i) IWFCMapGenerator {
	return &WFCMapGenerator{
		Size: size,
	}
}

func (wfc *WFCMapGenerator) Initialize(allPossibleModules []models.Module) bool {
	wfc.CurrentMap = make([][][][]models.Module, wfc.Size.X)
	for i := range wfc.CurrentMap {
		wfc.CurrentMap[i] = make([][][]models.Module, wfc.Size.Z)
		for j := range wfc.CurrentMap[i] {
			wfc.CurrentMap[i][j] = make([][]models.Module, wfc.Size.Y)
			for k := range wfc.CurrentMap[i][j] {
				modules := make([]models.Module, len(allPossibleModules))
				copy(modules, allPossibleModules)
				wfc.CurrentMap[i][j][k] = modules
			}
		}
	}
	return false
}

func (wfc WFCMapGenerator) Run(done chan bool) {
	for wfc.observe_and_collapse() {
		_ = wfc.propagate()
	}

	done <- true
}

func (wfc WFCMapGenerator) observe_and_collapse() bool {
	_, _ = wfc.calculate_lowest_entropy()
	return false
}

func (WFCMapGenerator) propagate() bool {
	return false
}

func (WFCMapGenerator) calculate_lowest_entropy() (*models.Vector3, bool) {
	return nil, false
}
