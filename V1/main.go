package main

import (
	"encoding/json"
	"os"

	"github.com/willcliffy/wfc_golang/generator"
	"github.com/willcliffy/wfc_golang/models"
)

var (
	EmptySlot            = ""
	SquareSlot           = "square"
	BottomRectangleSlot  = "rect"
	PositiveRampUpSlot   = "rampUp"
	PositiveRampDownSlot = "rampDown"
	NegativeRampUpSlot   = PositiveRampDownSlot
	NegativeRampDownSlot = PositiveRampUpSlot
)

func main() {
	mapSize := models.Vector3i{X: 10, Y: 10, Z: 10}
	mapGenerator := generator.NewWFCMapGenerator(mapSize)

	var allPossibleModules []models.Module

	// TODO - automate loading and analyzing cube meshes
	emptyModule := models.NewModule("empty", models.ModuleSlots{
		PositiveX: EmptySlot,
		PositiveY: EmptySlot,
		PositiveZ: EmptySlot,
		NegativeX: EmptySlot,
		NegativeY: EmptySlot,
		NegativeZ: EmptySlot,
	}, models.RotationY_000)
	allPossibleModules = append(allPossibleModules, emptyModule)

	cubeModule := models.NewModule("cube.glb", models.ModuleSlots{
		PositiveX: SquareSlot,
		PositiveY: SquareSlot,
		PositiveZ: SquareSlot,
		NegativeX: SquareSlot,
		NegativeY: SquareSlot,
		NegativeZ: SquareSlot,
	}, models.RotationY_000)
	allPossibleModules = append(allPossibleModules, cubeModule)

	flatModule := models.NewModule("flat.glb", models.ModuleSlots{
		PositiveX: BottomRectangleSlot,
		PositiveY: EmptySlot,
		PositiveZ: BottomRectangleSlot,
		NegativeX: BottomRectangleSlot,
		NegativeY: SquareSlot,
		NegativeZ: BottomRectangleSlot,
	}, models.RotationY_000)
	allPossibleModules = append(allPossibleModules, flatModule)

	cornerModuleSet := models.NewModuleSet("corner.glb", models.ModuleSlots{
		PositiveX: PositiveRampUpSlot,
		PositiveY: EmptySlot,
		PositiveZ: BottomRectangleSlot,
		NegativeX: BottomRectangleSlot,
		NegativeY: SquareSlot,
		NegativeZ: NegativeRampDownSlot,
	})
	allPossibleModules = append(allPossibleModules, cornerModuleSet...)

	rampModuleSet := models.NewModuleSet("ramp.glb", models.ModuleSlots{
		PositiveX: SquareSlot,
		PositiveY: EmptySlot,
		PositiveZ: PositiveRampUpSlot,
		NegativeX: BottomRectangleSlot,
		NegativeY: SquareSlot,
		NegativeZ: NegativeRampDownSlot,
	})
	allPossibleModules = append(allPossibleModules, rampModuleSet...)

	ok := mapGenerator.Initialize(allPossibleModules)
	if !ok {
		panic("failed to initialize map generator")
	}

	done := make(chan bool)
	go mapGenerator.Run(done)
	<-done

	finalMap := mapGenerator.GetFinalMap()
	jsonData, err := json.Marshal(finalMap)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("modules.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}
}
