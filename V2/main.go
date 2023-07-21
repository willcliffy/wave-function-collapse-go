package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	SIZE = Vector3i{25, 3, 25}
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
	seed := int64(0)
	wfc.Initialize(seed, SIZE, prototypes)

	doneChan := make(chan bool)
	go wfc.Run(doneChan)
	<-doneChan

	wfcMap := wfc.GetFinalMap(true)

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
