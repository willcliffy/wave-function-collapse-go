extends Control

var prototype_instance = preload("res://scenes/components/PrototypeButton.tscn")
var mesh_instance = preload("res://wfc_modules.glb")

const NUM_BUTTONS_IN_ROW = 4

var current_rows = 0
var current_buttons_in_row = 0


func _enter_tree():
	var x = 0
	for proto in Main.prototype_data_v1:
		var mesh_name = Main.prototype_data_v1[proto]["mesh_name"]
		if mesh_name == "-1":
			continue
		
		spawn_mesh(proto, Main.prototype_data_v1[proto], Vector3(x, 0, 0))
		x += 10


func spawn_mesh(proto_name: String, proto_data: Dictionary, mesh_position: Vector3):
	var proto = prototype_instance.instantiate()
	var mesh = mesh_instance.instantiate().get_node(proto_data["mesh_name"])

	proto.set_mesh(proto_data, mesh.mesh, mesh_position)
	proto.set_label(proto_name)
	proto.pressed.connect(
		func():
			Main.selected_prototype = proto_data
			Main.change_scene("res://scenes/screens/PrototypePreviewScreen.tscn")
	)

	get_node("Scroll/VBox/"+str(current_rows)).add_child(proto)
	current_buttons_in_row += 1

	if current_buttons_in_row == NUM_BUTTONS_IN_ROW:
		current_rows += 1
		current_buttons_in_row = 0
		var row = HBoxContainer.new()
		row.name = str(current_rows)
		$Scroll/VBox.add_child(row)
