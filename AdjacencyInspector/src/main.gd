extends CanvasLayer

@onready var meshes = preload("res://wfc_modules.glb")

const MAP_FILE_NAME = "prototype_data.json"

const STARTING_SCENE = "res://scenes/screens/AdjacencyInspectorScreen.tscn"

var selected_prototype
var prototype_data_v1

var scene_stack = []


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

	get_tree().change_scene_to_file(STARTING_SCENE)
	scene_stack.append(STARTING_SCENE)
	$AnimationPlayer.play_backwards("dissove")


func push_scene(target: String) -> void:
	print("Scene changed to: ", target)
	scene_stack.append(target)
	$AnimationPlayer.play("dissove")
	await $AnimationPlayer.animation_finished
	get_tree().change_scene_to_file(target)
	$AnimationPlayer.play_backwards("dissove")

func pop_scene() -> void:
	if len(scene_stack) <= 1:
		return

	scene_stack.pop_back()
	print("Scene changed to: ", scene_stack[0])

	$AnimationPlayer.play("dissove")
	await $AnimationPlayer.animation_finished
	get_tree().change_scene_to_file(scene_stack[0])
	$AnimationPlayer.play_backwards("dissove")


func _input(event):
	if event.is_action_pressed("back"):
		pop_scene()
