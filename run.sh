#!/bin/bash

# Run `Blender/generate_adjacency_rules.py` in Blender first

cp Blender/wfc_modules.glb Godot/wfc_modules.glb

cd Golang
go run . --input ../Blender/prototype_data.json --output ../Godot/map.json
