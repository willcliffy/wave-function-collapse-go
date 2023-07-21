package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	SIZE = Vector3i{12, 2, 12}
)

func main() {
	prototypeFileFlag := flag.String("input", "prototype_data.json", "Input file containing module (prototype) data")
	outputFileFlag := flag.String("output", "map.json", "Output file path/name")
	flag.Parse()

	if prototypeFileFlag == nil || *prototypeFileFlag == "" {
		panic("No prototype file provided")
	}

	prototypes, err := loadPrototypeDataFromFile(*prototypeFileFlag)
	if err != nil {
		panic(err)
	}

	wfc := NewWFCModel()
	wfc.Initialize(SIZE, prototypes)

	// TODO
	// apply_custom_constraints()

	doneChan := make(chan bool)
	go wfc.Run(doneChan)
	<-doneChan

	wfcMap := wfc.GetFinalMap()

	err = saveMapDataToFile(wfcMap, *outputFileFlag)
	if err != nil {
		panic(err)
	}
}

func loadPrototypeDataFromFile(filename string) (map[string]WFCPrototype, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var prototypes map[string]WFCPrototype
	err = json.Unmarshal(fileBytes, &prototypes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return prototypes, nil
}

func saveMapDataToFile(wfcMap interface{}, filePath string) error {
	prototypeJSON, err := json.Marshal(wfcMap)
	if err != nil {
		return fmt.Errorf("failed to marshal prototype to JSON: %v", err)
	}

	err = os.WriteFile(filePath, prototypeJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %v", err)
	}

	return nil
}

// func apply_custom_constraints():
// 	# This function isn't covered in the video but what we do here is basically
// 	# go over the wavefunction and remove certain modules from specific places
// 	# for example in my Blender scene I've marked all of the beach tiles with
// 	# an attribute called "constrain_to" with the value "bot". This is recalled
// 	# in this function, and all tiles with this attribute and value are removed
// 	# from cells that are not at the bottom i.e., if y > 0: constrain.
// 	var _add_to_stack = []

// 	for z in range(size.z):
// 		for y in range(size.y):
// 			for x in range(size.x):
// 				coords = Vector3(x, y, z)
// 				var protos = wfc.get_possibilities(coords)
// 				if y == size.y - 1:  # constrain top layer to not contain any uncapped prototypes
// 					for proto in protos.duplicate():
// 						var neighs  = protos[proto][WFC3D_Model.NEIGHBOURS][WFC3D_Model.pZ]
// 						if not "p-1" in neighs:
// 							protos.erase(proto)
// 							if not coords in wfc.stack:
// 								wfc.stack.append(coords)
// 				if y > 0:  # everything other than the bottom
// 					for proto in protos.duplicate():
// 						var custom_constraint = protos[proto][WFC3D_Model.CONSTRAIN_TO]
// 						if custom_constraint == WFC3D_Model.CONSTRAINT_BOTTOM:
// 							protos.erase(proto)
// 							if not coords in wfc.stack:
// 								wfc.stack.append(coords)
// 				if y < size.y - 1:  # everything other than the top
// 					for proto in protos.duplicate():
// 						var custom_constraint = protos[proto][WFC3D_Model.CONSTRAIN_TO]
// 						if custom_constraint == WFC3D_Model.CONSTRAINT_TOP:
// 							protos.erase(proto)
// 							if not coords in wfc.stack:
// 								wfc.stack.append(coords)
// 				if y == 0:  # constrain bottom layer so we don't start with any top-cliff parts at the bottom
// 					for proto in protos.duplicate():
// 						var neighs  = protos[proto][WFC3D_Model.NEIGHBOURS][WFC3D_Model.nZ]
// 						var custom_constraint = protos[proto][WFC3D_Model.CONSTRAIN_FROM]
// 						if (not "p-1" in neighs) or (custom_constraint == WFC3D_Model.CONSTRAINT_BOTTOM):
// 							protos.erase(proto)
// 							if not coords in wfc.stack:
// 								wfc.stack.append(coords)
// 				if x == size.x - 1: # constrain +x
// 					for proto in protos.duplicate():
// 						var neighs  = protos[proto][WFC3D_Model.NEIGHBOURS][WFC3D_Model.pX]
// 						if not "p-1" in neighs:
// 							protos.erase(proto)
// 							if not coords in wfc.stack:
// 								wfc.stack.append(coords)
// 				if x == 0: # constrain -x
// 					for proto in protos.duplicate():
// 						var neighs  = protos[proto][WFC3D_Model.NEIGHBOURS][WFC3D_Model.nX]
// 						if not "p-1" in neighs:
// 							protos.erase(proto)
// 							if not coords in wfc.stack:
// 								wfc.stack.append(coords)
// 				if z == size.z - 1: # constrain +z
// 					for proto in protos.duplicate():
// 						var neighs  = protos[proto][WFC3D_Model.NEIGHBOURS][WFC3D_Model.nY]
// 						if not "p-1" in neighs:
// 							protos.erase(proto)
// 							if not coords in wfc.stack:
// 								wfc.stack.append(coords)
// 				if z == 0: # constrain -z
// 					for proto in protos.duplicate():
// 						var neighs  = protos[proto][WFC3D_Model.NEIGHBOURS][WFC3D_Model.pY]
// 						if not "p-1" in neighs:
// 							protos.erase(proto)
// 							if not coords in wfc.stack:
// 								wfc.stack.append(coords)
