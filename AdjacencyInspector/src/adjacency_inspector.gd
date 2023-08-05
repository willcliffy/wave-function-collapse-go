extends Control

var prototype_instance = preload("res://scenes/components/PrototypeButton.tscn")


const NUM_BUTTONS_IN_ROW = 4

var current_rows = 0
var current_buttons_in_row = 0

var free_queue = []


func _enter_tree():
	var x = 0
	for proto_name in Main.prototype_data_v1:
		var proto_datum =  Main.prototype_data_v1[proto_name]
		var mesh_name = proto_datum["mesh_name"]
		if mesh_name == "-1":
			continue
		
		var mesh_position = Vector3(x, 0, 0)
		var mesh_rotation = Vector3(0, proto_datum["mesh_rotation"] * PI/2, 0)

		spawn_mesh(proto_name, proto_datum, mesh_position, mesh_rotation)
		x += 10

func _exit_tree():
	for _i in range(len(free_queue)):
		free_queue.pop_back().queue_free()


func spawn_mesh(proto_name: String, proto_data: Dictionary, mesh_position: Vector3, mesh_rotation: Vector3):
	var proto = prototype_instance.instantiate()
	var mesh = Main.meshes.instantiate().get_node(proto_data["mesh_name"])

	proto.set_preview(mesh.mesh, mesh_position, mesh_rotation)
	proto.set_label(proto_name)
	proto.get_node("Prototype").pressed.connect(
		func():
			Main.selected_prototype = proto_name
			Main.push_scene("res://scenes/screens/PrototypePreviewScreen.tscn")
	)

	get_node("Scroll/VBox/"+str(current_rows)).add_child(proto)
	free_queue.append(proto)
	current_buttons_in_row += 1

	if current_buttons_in_row == NUM_BUTTONS_IN_ROW:
		current_rows += 1
		current_buttons_in_row = 0
		var row = HBoxContainer.new()
		row.name = str(current_rows)
		$Scroll/VBox.add_child(row)
