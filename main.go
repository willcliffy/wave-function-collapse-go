package main

import (
	"github.com/willcliffy/wfc_golang/generator"
	"github.com/willcliffy/wfc_golang/models"
)

func main() {
	mapSize := models.Vector3i{X: 10, Y: 10, Z: 10}
	mapGenerator := generator.NewWFCMapGenerator(mapSize)

	var allPossibleModules []models.Module

	// TODO - load modules here
	flatModuleSet := models.NewModuleSet("", models.ModuleSlots{})
	allPossibleModules = append(allPossibleModules, flatModuleSet...)

	rampModuleSet := models.NewModuleSet("", models.ModuleSlots{})
	allPossibleModules = append(allPossibleModules, rampModuleSet...)

	cubeModuleSet := models.NewModuleSet("", models.ModuleSlots{})
	allPossibleModules = append(allPossibleModules, cubeModuleSet...)

	cornerModuleSet := models.NewModuleSet("", models.ModuleSlots{})
	allPossibleModules = append(allPossibleModules, cornerModuleSet...)

	ok := mapGenerator.Initialize(allPossibleModules)
	if !ok {
		panic("failed to initialize map generator")
	}

	done := make(chan bool)
	mapGenerator.Run(done)
	<-done
}
