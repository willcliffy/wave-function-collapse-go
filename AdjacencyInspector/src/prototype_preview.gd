extends Control

var proto_preview = preload("res://scenes/components/PrototypePreview.tscn")
var mesh_instance = preload("res://wfc_modules.glb")


func _enter_tree():
	var proto_mesh = mesh_instance.instantiate().get_node(Main.selected_prototype["mesh_name"])
	var proto = proto_preview.instantiate()
	proto.name = Main.selected_prototype["mesh_name"]
	var proto_position = Vector3(1, 1, 1)
	var proto_rotation = Vector3(0, Main.selected_prototype["mesh_rotation"], 0)
	proto.set_mesh(proto_mesh.mesh, proto_position, proto_rotation)
	$Container/Preview/SubViewportContainer/SubViewport.add_child(proto)
