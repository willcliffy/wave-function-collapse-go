extends Node3D

var mesh_instance


func set_mesh(new_mesh: ArrayMesh, mesh_position: Vector3, mesh_rotation: Vector3):
	mesh_instance = MeshInstance3D.new()
	mesh_instance.mesh = new_mesh
	mesh_instance.rotation = mesh_rotation
	$Node.add_child(mesh_instance)
	$Node.position = mesh_position
	print("node at pos ", mesh_position)


func _process(delta):
	if mesh_instance != null:
		mesh_instance.rotate_y(delta)
