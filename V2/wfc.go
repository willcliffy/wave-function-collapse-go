package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
)

const (
	MESH_NAME         = "mesh_name"
	MESH_ROT          = "mesh_rotation"
	NEIGHBOURS        = "valid_neighbours"
	CONSTRAIN_TO      = "constrain_to"
	CONSTRAIN_FROM    = "constrain_from"
	CONSTRAINT_BOTTOM = "bot"
	CONSTRAINT_TOP    = "top"
	WEIGHT            = "weight"

	P_X = 0
	P_Y = 1
	N_X = 2
	N_Y = 3
	P_Z = 4
	N_Z = 5
)

var (
	DIRECTION_TO_INDEX = map[Vector3i]int{
		V3i_RIGHT:   0,
		V3i_FORWARD: 1, // should be 3?
		V3i_LEFT:    2,
		V3i_BACK:    3, // should be 1?
		V3i_UP:      4,
		V3i_DOWN:    5,
	}
)

type WFCModel interface {
	Initialize(newSize Vector3i, allPrototypes map[string]WFCPrototype)
	Run(chan bool)
	GetFinalMap() *WFCMapLinear
}

type WFC struct {
	waveFunction [][][]map[string]WFCPrototype
	finalMap     WFCMapLinear
	size         Vector3i
}

func NewWFCModel() WFCModel {
	return &WFC{}
}

func deepCopy(src, dst interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buf).Decode(dst)
}

func (wfc *WFC) Initialize(newSize Vector3i, allPrototypes map[string]WFCPrototype) {
	wfc.size = newSize
	wfc.finalMap.Size = newSize
	for z := 0; z < wfc.size.Z; z++ {
		var ySlice [][]map[string]WFCPrototype
		for y := 0; y < wfc.size.Y; y++ {
			var xSlice []map[string]WFCPrototype
			for x := 0; x < wfc.size.X; x++ {
				var prototypeClone map[string]WFCPrototype
				err := deepCopy(allPrototypes, &prototypeClone)
				if err != nil {
					panic(err)
				}
				xSlice = append(xSlice, prototypeClone)
			}
			ySlice = append(ySlice, xSlice)
		}
		wfc.waveFunction = append(wfc.waveFunction, ySlice)
	}
}

func (wfc *WFC) Run(doneChan chan bool) {
	for !wfc.isCollapsed() {
		min_entropy_coords := wfc.getMinEntropyCoords()
		if min_entropy_coords == nil {
			fmt.Printf("min_entropy_coords are nil!")
			break
		}

		wfc.collapse(min_entropy_coords)
		fmt.Printf("tick\n")
		wfc.propagate(min_entropy_coords, false)
		fmt.Printf("tick\n")
	}

	doneChan <- true
}

func (wfc WFC) GetFinalMap() *WFCMapLinear {
	if !wfc.isCollapsed() {
		return nil
	}

	// result := WFCMap{
	// 	Size:       wfc.size,
	// 	Prototypes: make([][][]WFCPrototype, wfc.size.Z),
	// }

	// for z := 0; z < wfc.size.Z; z++ {
	// 	result.Prototypes[z] = make([][]WFCPrototype, wfc.size.Y)
	// 	for y := 0; y < wfc.size.Y; y++ {
	// 		result.Prototypes[z][y] = make([]WFCPrototype, wfc.size.X)
	// 		for x := 0; x < wfc.size.Z; x++ {
	// 			prototypes := wfc.waveFunction[z][y][x]

	// 			// only grab collapsed cells (?)
	// 			if len(prototypes) > 1 {
	// 				fmt.Printf("[WARN] Uncollapsed cell in GetFinalMap")
	// 				continue
	// 			}

	// 			for _, prototype := range prototypes {
	// 				result.Prototypes[x][y][z] = prototype
	// 			}
	// 		}
	// 	}
	// }

	return &wfc.finalMap
}

func (wfc *WFC) propagate(coords *Vector3i, singleIteration bool) {
	var stack []Vector3i
	if coords != nil {
		stack = append(stack, *coords)
	}

	for len(stack) > 0 {
		stackLen := len(stack)
		currentCoords := stack[stackLen-1]
		stack = stack[:stackLen-1]

		for _, direction := range wfc.validDirections(&currentCoords) {
			otherCoords := currentCoords.Add(direction)
			otherProtos := wfc.waveFunction[otherCoords.Z][otherCoords.Y][otherCoords.X]

			if len(otherProtos) == 0 {
				continue
			}

			possibleNeighbors := wfc.getPossibleNeighbours(&currentCoords, direction)
			for otherProto := range otherProtos {
				if StringSliceContains(possibleNeighbors, otherProto) {
					continue
				}

				// Constrain
				delete(wfc.waveFunction[coords.Z][coords.Y][coords.X], otherProto)

				if len(wfc.waveFunction[coords.Z][coords.Y][coords.X]) == 1 {
					for _, v := range wfc.waveFunction[coords.Z][coords.Y][coords.X] {
						wfc.finalMap.Prototypes = append(wfc.finalMap.Prototypes, v)
					}
				}

				if !Vector3iSliceContains(stack, otherCoords) {
					stack = append(stack, otherCoords)
					fmt.Printf("Stack doesnt contain %v - %v\n", otherCoords, stack)
				}
			}
		}

		if singleIteration {
			break
		}
	}
}

func (wfc *WFC) isCollapsed() bool {
	for _, z := range wfc.waveFunction {
		for _, y := range z {
			for _, x := range y {
				if len(x) > 1 {
					return false
				}
			}
		}
	}
	return true
}

func (wfc *WFC) getPossibleNeighbours(coords *Vector3i, dir Vector3i) []string {
	var validNeighbours []string
	prototypes := wfc.waveFunction[coords.Z][coords.Y][coords.X]
	dirIdx := DIRECTION_TO_INDEX[dir]
	for _, prototype := range prototypes {
		neighbours := prototype.ValidNeighbours[dirIdx]
		for _, neighbor := range neighbours {
			if !StringSliceContains(validNeighbours, neighbor) {
				validNeighbours = append(validNeighbours, neighbor)
			}
		}
	}
	return validNeighbours
}

func (wfc *WFC) collapse(coords *Vector3i) {
	possibleProtos := wfc.waveFunction[coords.Z][coords.Y][coords.X]
	protoName := wfc.weightedChoice(possibleProtos)
	wfc.waveFunction[coords.Z][coords.Y][coords.X] = map[string]WFCPrototype{protoName: possibleProtos[protoName]}
	wfc.finalMap.Prototypes = append(wfc.finalMap.Prototypes, possibleProtos[protoName])
}

func (wfc *WFC) weightedChoice(prototypes map[string]WFCPrototype) string {
	protoWeights := make(map[float64]string)

	for p, properties := range prototypes {
		w := float64(properties.Weight)
		w += rand.Float64()*(1.0-(-1.0)) + (-1.0)
		protoWeights[w] = p
	}

	weightList := make([]float64, 0, len(protoWeights))
	for weight := range protoWeights {
		weightList = append(weightList, weight)
	}
	sortFloat64Slice(weightList)

	protoName := protoWeights[weightList[len(weightList)-1]]

	return protoName
}

func (wfc *WFC) getMinEntropyCoords() *Vector3i {
	minEntropy := math.MaxInt
	var coords *Vector3i
	for z := 0; z < wfc.size.Z; z++ {
		for y := 0; y < wfc.size.Y; y++ {
			for x := 0; x < wfc.size.Z; x++ {
				entropy := len(wfc.waveFunction[z][y][x])
				if entropy > 1 {
					entropy += int(rand.Float64()*0.2 - 0.1) // Add random float between -0.1 and 0.1

					if entropy < minEntropy {
						minEntropy = entropy
						coords = &Vector3i{X: x, Y: y, Z: z}
					}
				}
			}
		}
	}
	return coords
}

func (wfc *WFC) validDirections(coords *Vector3i) []Vector3i {
	var dirs []Vector3i

	if coords.X > 0 {
		dirs = append(dirs, V3i_LEFT)
	}
	if coords.X < wfc.size.X-1 {
		dirs = append(dirs, V3i_RIGHT)
	}
	if coords.Y > 0 {
		dirs = append(dirs, V3i_DOWN)
	}
	if coords.Y < wfc.size.Y-1 {
		dirs = append(dirs, V3i_UP)
	}
	if coords.Z > 0 {
		dirs = append(dirs, V3i_FORWARD)
	}
	if coords.Z < wfc.size.Z-1 {
		dirs = append(dirs, V3i_BACK)
	}

	return dirs
}
