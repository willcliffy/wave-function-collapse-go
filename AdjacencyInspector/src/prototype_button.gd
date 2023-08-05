extends Control

var proto_preview = preload("res://scenes/components/PrototypePreview.tscn")

var protos = []


func set_preview(new_mesh: ArrayMesh, mesh_position: Vector3, mesh_rotation: Vector3):
	for i in range(len(protos)):
		var proto = protos.pop_front()
		proto.queue_free()
		
	var proto = proto_preview.instantiate()
	proto.set_mesh(new_mesh, mesh_position, mesh_rotation)
	$Prototype/VPContainer/VP.add_child(proto)
	protos.append(proto)


func set_custom_scale(custom_scale):
	$Prototype.scale = Vector2.ONE * custom_scale
	$Label.scale *=  Vector2.ONE * custom_scale
	custom_minimum_size *=  Vector2.ONE * custom_scale


func set_label(label_text):
	$Label.text = label_text


func set_selected():
	pass
