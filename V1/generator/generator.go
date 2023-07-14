package generator

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/willcliffy/wfc_golang/models"
)

type IWFCMapGenerator interface {
	Initialize(allPossibleModules []models.Module) bool
	Run(done chan bool)
	GetFinalMap() models.WFCMap
	observe() (cell *models.Vector3i, entropy int, done bool)
	collapse(cell *models.Vector3i, entropy int) (*models.Module, bool)
	propagate(cell *models.Vector3i, modules []*models.Module)
}

type WFCMapGenerator struct {
	randIntGen *rand.Rand

	Size models.Vector3i

	// x, z, y -> list of potential modules
	CurrentMap [][][][]*models.Module
}

func NewWFCMapGenerator(size models.Vector3i) IWFCMapGenerator {
	return &WFCMapGenerator{
		randIntGen: rand.New(rand.NewSource(time.Now().UnixNano())),
		Size:       size,
	}
}

func (wfc WFCMapGenerator) GetFinalMap() models.WFCMap {
	var finalMap models.WFCMap
	for i := range wfc.CurrentMap {
		for j := range wfc.CurrentMap[i] {
			for _, cell := range wfc.CurrentMap[i][j] {
				if len(cell) > 1 {
					fmt.Printf("Non-collapsed cell in GetFinalMap, ignoring\n")
				} else if len(cell) == 0 {
					// fmt.Printf("Unassigned cell in GetFinalMap, ignoring\n")
					continue
				} else if cell[0].Filename == "empty" {
					continue
				}

				finalMap = append(finalMap, cell[0].ToWFCMapCell())
			}
		}
	}
	return finalMap
}

func (wfc *WFCMapGenerator) Initialize(allPossibleModules []models.Module) bool {
	wfc.CurrentMap = make([][][][]*models.Module, wfc.Size.X)
	for i := range wfc.CurrentMap {
		wfc.CurrentMap[i] = make([][][]*models.Module, wfc.Size.Z)
		for j := range wfc.CurrentMap[i] {
			wfc.CurrentMap[i][j] = make([][]*models.Module, wfc.Size.Y)
			for k := range wfc.CurrentMap[i][j] {
				wfc.CurrentMap[i][j][k] = make([]*models.Module, len(allPossibleModules))
				modules := make([]models.Module, len(allPossibleModules))
				copy(modules, allPossibleModules)
				for l, module := range modules {
					module.Position.X = i
					module.Position.Z = j
					module.Position.Y = k
					wfc.CurrentMap[i][j][k][l] = module.Copy()
				}
			}
		}
	}

	// var cube models.Module
	// for _, mod := range allPossibleModules {
	// 	if mod.Filename == "cube.glb" {
	// 		cube = mod
	// 		break
	// 	}
	// }
	// cell := models.Vector3i{X: wfc.Size.X / 2, Y: wfc.Size.Y / 2, Z: wfc.Size.Z / 2}
	// wfc.CurrentMap[cell.X][cell.Z][cell.Y] = []*models.Module{&cube}
	// wfc.propagate(&cell, []*models.Module{&cube})
	// for i := range wfc.CurrentMap {
	// 	fmt.Println(wfc.CurrentMap[i])
	// }

	return true
}

func (wfc WFCMapGenerator) Run(doneChan chan bool) {
	s := time.Now()
	cell, entropy, done := wfc.observe()
	for !done {
		module, _ := wfc.collapse(cell, entropy)
		wfc.propagate(cell, []*models.Module{module})
		cell, entropy, done = wfc.observe()
	}

	elapsed := time.Since(s)
	fmt.Printf("Done in %.2d:%.2d.%.4d\n", int(elapsed.Minutes()), int(elapsed.Seconds()), elapsed.Milliseconds())
	doneChan <- true
}

func (wfc WFCMapGenerator) observe() (cell *models.Vector3i, entropy int, done bool) {
	entropy = math.MaxInt

	for i := range wfc.CurrentMap {
		for j := range wfc.CurrentMap[i] {
			for k := range wfc.CurrentMap[i][j] {
				cellEntropy := len(wfc.CurrentMap[i][j][k])
				if cellEntropy >= entropy {
					continue
				}

				if cellEntropy == 1 {
					// May want to change this, but for now, a cell entropy of 1 means that this
					// cell has been collapsed, either explicitly or implicitly
					continue
				}

				if cellEntropy == 0 {
					// This indicates a failure case; there are no modules which fit in this cell
					// A few things can be done here, but for now, we silently continue
					// - Backtrack
					// - Insert new tile which handles this case gracefully?
					// - Error out
					// - skip cell and continue (current option)
					continue
				}

				entropy = cellEntropy
				cell = &models.Vector3i{X: i, Z: j, Y: k}
			}
		}
	}

	if entropy == math.MaxInt {
		done = true
	}

	return
}

func (wfc *WFCMapGenerator) collapse(cell *models.Vector3i, entropy int) (*models.Module, bool) {
	chosenModuleInd := wfc.randIntGen.Intn(entropy)
	chosenModule := wfc.CurrentMap[cell.X][cell.Z][cell.Y][chosenModuleInd]
	wfc.CurrentMap[cell.X][cell.Z][cell.Y] = []*models.Module{chosenModule}

	return chosenModule, true
}

func (wfc *WFCMapGenerator) propagate(cell *models.Vector3i, modules []*models.Module) {
	changed := make(map[models.Vector3i][]*models.Module)

	for i := cell.X - 1; i <= cell.X+1; i++ {
		if i < 0 || i >= len(wfc.CurrentMap) {
			continue
		}
		for j := cell.Z - 1; j <= cell.Z+1; j++ {
			if j < 0 || j >= len(wfc.CurrentMap[i]) {
				continue
			}
			for k := cell.Y - 1; k <= cell.Y+1; k++ {
				if k < 0 || k >= len(wfc.CurrentMap[i][j]) {
					continue
				}

				if i == cell.X && j == cell.Z && k == cell.Y {
					continue
				}

				// For now, we are only concerned with the tiles directly adjacent to the cell
				if cell.Distance(models.Vector3i{X: i, Z: j, Y: k}) > 1 {
					continue
				}

				cellChanged := false
				var newCell []*models.Module
				for _, possibleModule := range wfc.CurrentMap[i][j][k] {
					compatibleModuleFound := false
					for _, module := range modules {
						if module.CompatibleWith(possibleModule) {
							compatibleModuleFound = true
							break
						} else {
							if possibleModule.Filename == "cube.glb" && module.Filename == "ramp.glb" {
								fmt.Printf("cube at %v not compat with ramp: %+v, %v\n", possibleModule.Position, module.Position, module.Rotation)
							}
						}
					}

					if compatibleModuleFound {

						newCell = append(newCell, possibleModule)
					} else {
						cellChanged = true
					}
				}

				if cellChanged {
					wfc.CurrentMap[i][j][k] = newCell
					if len(newCell) > 0 {
						changed[models.Vector3i{X: i, Z: j, Y: k}] = wfc.CurrentMap[i][j][k]
					}
				}
			}
		}
	}

	for cell, modules := range changed {
		wfc.propagate(&cell, modules)
	}
}
