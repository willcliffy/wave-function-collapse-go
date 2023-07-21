extends Node3D

var mesh : ArrayMesh: set = set_mesh
var prototype : Dictionary
var debug_text

var text


func set_mesh(new_mesh):
	mesh = new_mesh
	$mesh_instance.mesh = mesh


func _on_col_area_mouse_entered():
	if debug_text:
		debug_text.text = str(prototype)
	$Highlight.visible = true


func _on_col_area_mouse_exited():
	if debug_text and debug_text.text:
		debug_text.text = ""
	$Highlight.visible = false
