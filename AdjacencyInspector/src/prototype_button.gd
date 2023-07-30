extends Button

var proto_preview = preload("res://scenes/components/PrototypePreview.tscn")


func set_mesh(proto_data: Dictionary, new_mesh: ArrayMesh, mesh_position: Vector3):
	var proto = proto_preview.instantiate()
	var mesh_rotation = Vector3(0, proto_data["mesh_rotation"] * PI/2, 0)
	proto.set_mesh(new_mesh, mesh_position, mesh_rotation)
	$VPContainer/VP.add_child(proto)
	print(mesh_rotation)


func set_label(label_text):
	$Label.text = label_text
