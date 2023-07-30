extends CanvasLayer

const MAP_FILE_NAME = "prototype_data.json"

var selected_prototype
var prototype_data_v1


func _enter_tree():
	if not FileAccess.file_exists(MAP_FILE_NAME):
		print("File not found.")
		return

	var file = FileAccess.open(MAP_FILE_NAME, FileAccess.READ)
	var json_text = file.get_as_text()
	file.close()

	var json_dict = JSON.parse_string(json_text)
	if typeof(json_dict) != TYPE_DICTIONARY:
		print("Failed to parse JSON.")
		return

	prototype_data_v1 = json_dict

	get_tree().change_scene_to_file("res://scenes/screens/AdjacencyInspectorScreen.tscn")
	$AnimationPlayer.play_backwards("dissove")


func change_scene(target: String) -> void:
	print("Scene changed to: ", target)
	$AnimationPlayer.play("dissove")
	await $AnimationPlayer.animation_finished
	get_tree().change_scene_to_file(target)
	$AnimationPlayer.play_backwards("dissove")


func _input(event):
	if event.is_action_pressed("back"):
		print("back nyi")
