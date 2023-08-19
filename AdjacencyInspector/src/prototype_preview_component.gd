extends Node3D

var neighbors = {}


func set_mesh(new_mesh: ArrayMesh, mesh_position: Vector3, mesh_rotation: Vector3):
	$Space/Node3D.rotation = Vector3.ZERO
	var mesh_instance = MeshInstance3D.new()
	mesh_instance.mesh = new_mesh
	mesh_instance.rotation = mesh_rotation
	$Space/Node3D.add_child(mesh_instance)
	position = mesh_position


func set_neighbor(direction, new_mesh, mesh_rotation):
	if direction in neighbors:
		neighbors[direction].queue_free()
		neighbors.erase(direction)
	if new_mesh == null:
		return
	var mesh_instance = MeshInstance3D.new()
	mesh_instance.mesh = new_mesh
	mesh_instance.rotation = mesh_rotation
	$Space/Node3D.add_child(mesh_instance)
	mesh_instance.position += direction
	neighbors[direction] = mesh_instance


func _process(delta):
	$Space.rotate_y(delta)
