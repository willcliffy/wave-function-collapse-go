package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"time"
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
	Initialize(seed int64, newSize Vector3i, allPrototypes map[string]WFCPrototype)
	Run(done chan bool)
	GetFinalMap(onlyCollapsed bool) *WFCMap
}

type WFC struct {
	waveFunction [][][]map[string]WFCPrototype
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

func (wfc *WFC) Initialize(seed int64, newSize Vector3i, allPrototypes map[string]WFCPrototype) {
	rand.Seed(seed)

	wfc.size = newSize
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

	stack := wfc.applyCustomConstraints()
	wfc.propagate(stack, false)
}

func (wfc *WFC) applyCustomConstraints() []Vector3i {
	var stack []Vector3i

	for z := 0; z < wfc.size.Z; z++ {
		for y := 0; y < wfc.size.Y; y++ {
			for x := 0; x < wfc.size.X; x++ {
				coords := Vector3i{x, y, z}
				protos := wfc.waveFunction[z][y][x]

				// constrain top layer to not contain any uncapped prototypes
				if y == wfc.size.Y-1 {
					for proto := range duplicateMap(protos) {
						neighs := protos[proto].ValidNeighbours[P_Z]
						if !StringSliceContains(neighs, "p-1") {
							// fmt.Printf("Constraining %v from being %v since it's uncapped\n", coords, proto)
							delete(protos, proto)
							if !Vector3iSliceContains(stack, coords) {
								stack = append(stack, coords)
							}
						}
					}
				}

				// everything other than the bottom
				if y > 0 {
					for proto := range duplicateMap(protos) {
						customConstraint := protos[proto].ConstrainTo
						if customConstraint == CONSTRAINT_BOTTOM {
							// fmt.Printf("Constraining %v from being %v since it has constraint BOT\n", coords, proto)
							delete(protos, proto)
							if !Vector3iSliceContains(stack, coords) {
								stack = append(stack, coords)
							}
						}
					}
				}

				// everything other than the top
				if y < wfc.size.Y-1 {
					for proto := range duplicateMap(protos) {
						customConstraint := protos[proto].ConstrainTo
						if customConstraint == CONSTRAINT_TOP {
							// fmt.Printf("Constraining %v from being %v since it has constraint TOP\n", coords, proto)
							delete(protos, proto)
							if !Vector3iSliceContains(stack, coords) {
								stack = append(stack, coords)
							}
						}
					}
				}

				// constrain bottom layer so we don't start with any top-cliff parts at the bottom
				if y == 0 {
					for proto := range duplicateMap(protos) {
						neighs := protos[proto].ValidNeighbours[N_Z]
						customConstraint := protos[proto].ConstrainFrom
						if !StringSliceContains(neighs, "p-1") || customConstraint == CONSTRAINT_BOTTOM {
							// fmt.Printf("Constraining %v from being %v since it's a top-cliff part\n", coords, proto)
							delete(protos, proto)
							if !Vector3iSliceContains(stack, coords) {
								stack = append(stack, coords)
							}
						}
					}
				}

				wfc.waveFunction[z][y][x] = protos
			}
		}
	}

	return stack
}

func (wfc *WFC) Run(doneChan chan bool) {
	s := time.Now()
	for !wfc.isCollapsed() {
		min_entropy_coords := wfc.getMinEntropyCoords()
		if min_entropy_coords == nil {
			fmt.Printf("min_entropy_coords are nil!")
			break
		}

		wfc.collapse(min_entropy_coords)
		wfc.propagateCoord(min_entropy_coords, false)
		// fmt.Printf("tick\n")
	}

	fmt.Printf("WFC Collapsed in %v\n", time.Since(s))

	doneChan <- true
}

func (wfc WFC) GetFinalMap(onlyCollapsed bool) *WFCMap {
	if !wfc.isCollapsed() {
		return nil
	}

	var wfcMap WFCMap
	wfcMap.Size = wfc.size

	wfcMap.Prototypes = make([][][]WFCPrototypeFinalized, wfc.size.Z)
	for z := range wfc.waveFunction {
		wfcMap.Prototypes[z] = make([][]WFCPrototypeFinalized, wfc.size.Y)
		for y := range wfc.waveFunction[z] {
			wfcMap.Prototypes[z][y] = make([]WFCPrototypeFinalized, wfc.size.X)
			for x, cell := range wfc.waveFunction[z][y] {
				if len(cell) > 1 {
					fmt.Printf("[WARN] Uncollapsed cell in GetFinalMap: (%d, %d, %d) - %v\n", x, y, z, cell)
					if onlyCollapsed {
						continue
					}
				} else if len(cell) == 0 {
					fmt.Printf("[WARN] Nil'd cell in GetFinalMap: (%d, %d, %d) - %v\n", x, y, z, cell)
				}

				for _, proto := range cell {
					wfcMap.Prototypes[z][y][x] = *proto.Finalize(&Vector3i{X: x, Y: y, Z: z})
				}
			}
		}
	}

	return &wfcMap
}

func (wfc *WFC) propagateCoord(coords *Vector3i, singleIteration bool) {
	var stack []Vector3i
	if coords != nil {
		stack = append(stack, *coords)
	}
	wfc.propagate(stack, singleIteration)
}

func (wfc *WFC) propagate(stack []Vector3i, singleIteration bool) {
	for len(stack) > 0 {
		stackLen := len(stack)
		coords := stack[stackLen-1]
		stack = stack[:stackLen-1]

		for _, direction := range wfc.validDirections(&coords) {
			otherCoords := coords.Add(direction)

			otherProtos := wfc.waveFunction[otherCoords.Z][otherCoords.Y][otherCoords.X]
			if len(otherProtos) == 0 {
				continue
			}

			possibleNeighbors := wfc.getPossibleNeighbours(&coords, direction)
			for otherProto := range otherProtos {
				if StringSliceContains(possibleNeighbors, otherProto) {
					continue
				}

				// Constrain
				delete(wfc.waveFunction[otherCoords.Z][otherCoords.Y][otherCoords.X], otherProto)

				if !Vector3iSliceContains(stack, otherCoords) {
					// fmt.Printf("Stack doesnt contain  %v at %v\n", otherCoords, coords)
					// fmt.Printf("\tStack:              %v\n", stack)
					// fmt.Printf("\tPossible Neighbors: %v\n", possibleNeighbors)
					// fmt.Printf("\tPossibilities:      %v\n", wfc.waveFunction[coords.Z][coords.Y][coords.X])
					// fmt.Printf("\tOther Proto:        %v\n", otherProto)
					stack = append(stack, otherCoords)
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
	proto := possibleProtos[protoName]
	wfc.waveFunction[coords.Z][coords.Y][coords.X] = map[string]WFCPrototype{protoName: proto}
}

func (wfc *WFC) weightedChoice(prototypes map[string]WFCPrototype) string {
	protoWeights := make(map[float64]string)

	for p, properties := range prototypes {
		w := float64(properties.Weight)
		w += 2.0*rand.Float64() - 1.0
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
