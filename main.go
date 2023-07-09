package main

import (
	"encoding/json"
	"os"

	"github.com/willcliffy/wfc_golang/generator"
	"github.com/willcliffy/wfc_golang/models"
)

var (
	EmptySlot           = ""
	SquareSlot          = "aaaa"
	BottomRectangleSlot = "bbbb"
	RampUpSlot          = "cccc"
	RampDownSlot        = "dddd"
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

	cornerModuleSet := models.NewModuleSet("cube.glb", models.ModuleSlots{
		PositiveX: SquareSlot,
		PositiveY: SquareSlot,
		PositiveZ: SquareSlot,
		NegativeX: SquareSlot,
		NegativeY: SquareSlot,
		NegativeZ: SquareSlot,
	})
	allPossibleModules = append(allPossibleModules, cornerModuleSet...)

	flatModuleSet := models.NewModuleSet("corner.glb", models.ModuleSlots{
		PositiveX: RampUpSlot,
		PositiveY: EmptySlot,
		PositiveZ: BottomRectangleSlot,
		NegativeX: BottomRectangleSlot,
		NegativeY: SquareSlot,
		NegativeZ: RampDownSlot,
	})
	allPossibleModules = append(allPossibleModules, flatModuleSet...)

	rampModuleSet := models.NewModuleSet("ramp.glb", models.ModuleSlots{
		PositiveX: SquareSlot,
		PositiveY: EmptySlot,
		PositiveZ: RampUpSlot,
		NegativeX: BottomRectangleSlot,
		NegativeY: SquareSlot,
		NegativeZ: RampDownSlot,
	})
	allPossibleModules = append(allPossibleModules, rampModuleSet...)

	cubeModuleSet := models.NewModuleSet("flat.glb", models.ModuleSlots{
		PositiveX: BottomRectangleSlot,
		PositiveY: EmptySlot,
		PositiveZ: BottomRectangleSlot,
		NegativeX: BottomRectangleSlot,
		NegativeY: SquareSlot,
		NegativeZ: BottomRectangleSlot,
	})
	allPossibleModules = append(allPossibleModules, cubeModuleSet...)

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

	err = os.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}
}
