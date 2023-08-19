#!/bin/bash

# Run `Blender/generate_adjacency_rules.py` in Blender first

cp Blender/wfc_modules.glb Godot/wfc_modules.glb

cp Blender/wfc_modules.glb AdjacencyInspector/wfc_modules.glb
cp Blender/prototype_data.json AdjacencyInspector/prototype_data.json

cd Golang
go run . --input ../Blender/prototype_data.json --output ../Godot/map.json
