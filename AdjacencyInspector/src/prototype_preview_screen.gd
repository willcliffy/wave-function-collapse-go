extends Control

@onready var prototype_instance = preload("res://scenes/components/PrototypeButton.tscn")
@onready var proto_preview = preload("res://scenes/components/PrototypePreview.tscn")
@onready var meshes = Main.meshes.instantiate()

# TODO: clean up
var BASE_POSITION = Vector3(100, 100, 100)
var BASE_POSITION_INCREMENT = Vector3(10, 0, 0)
var BASE_POSITION_INDEX = 0

var preview_3d

func _ready():
	var neighbors = Main.prototype_data_v1[Main.selected_prototype]["valid_neighbours"]
	var neighbor_ui_containers = [
		$Container/Neighbors/Top/Scroll,
		$Container/Neighbors/Bot/Scroll,
		$Container/Neighbors/PX/Scroll,
		$Container/Neighbors/PZ/Scroll,
		$Container/Neighbors/NX/Scroll,
		$Container/Neighbors/NZ/Scroll
	]

	var directions = [
		Vector3(0, 1, 0),
		Vector3(0, -1, 0),
		Vector3(1, 0, 0),
		Vector3(0, 0, 1),
		Vector3(-1, 0, 0),
		Vector3(0, 0, -1),
	]

	for i in range(len(neighbors)):
		load_direction(directions[i], neighbors[i], neighbor_ui_containers[i])

	var selected_prototype = Main.prototype_data_v1[Main.selected_prototype]

	var proto_mesh = meshes.get_node(selected_prototype["mesh_name"])
	preview_3d = proto_preview.instantiate()
	preview_3d.name = selected_prototype["mesh_name"]
	var proto_position = Vector3(1, 1, 1)
	var proto_rotation = Vector3(0, selected_prototype["mesh_rotation"], 0)
	preview_3d.set_mesh(proto_mesh.mesh, proto_position, proto_rotation)
	$Container/Preview/SubViewportContainer/SubViewport.add_child(preview_3d)


func load_direction(direction: Vector3, neighbors: Array, component: ScrollContainer):
	component.custom_minimum_size = Vector2(410, 600 * (1.0 / 0.7))
	for proto in neighbors:
		var proto_datum = Main.prototype_data_v1[proto]
		var mesh_name = proto_datum["mesh_name"]

		if mesh_name == "-1":
			continue

		var proto_button: Control = prototype_instance.instantiate()
		var mesh_position = BASE_POSITION + BASE_POSITION_INDEX * BASE_POSITION_INCREMENT
		var mesh_rotation = Vector3(0, proto_datum["mesh_rotation"] * PI/2, 0)
		var mesh = meshes.get_node(mesh_name).mesh
		proto_button.set_label(proto)
		proto_button.set_preview(mesh, mesh_position, mesh_rotation)
		proto_button.set_custom_scale(0.70)
		proto_button.get_node("Prototype").pressed.connect(
			func():
				preview_3d.set_neighbor(direction, mesh, mesh_rotation)
		)
		component.get_node("VB").add_child(proto_button)
		BASE_POSITION_INDEX += 1
