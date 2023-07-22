extends Node

@onready var module_scene = preload("res://Module.tscn").instantiate()
@onready var modules_scene = preload("res://wfc_modules.glb").instantiate()

const MAP_FILE_NAME = "map.json"


func spawn_mesh(mesh_name: String, position: Vector3, rotation: int):
	var inst = module_scene.duplicate()
	add_child(inst)
	# meshes.append(inst) # TODO - keep track and free
	inst.position = position
	inst.rotate_y((PI/2) * rotation)
	var mesh = modules_scene.get_node(mesh_name).mesh.duplicate()
	inst.mesh = mesh


func _ready():
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

	var size = Vector3(json_dict["Size"]["X"], json_dict["Size"]["Y"], json_dict["Size"]["Z"])
	print("Size: ", size)
	$CameraBase.position = size / 2.0

	var prototypes = json_dict["Prototypes"]
	for z in range(size.z):
		for y in range(size.y):
			for x in range(size.x):
				var prototype = prototypes[z][y][x]

				var mesh_name = prototype["mesh_name"]
				if mesh_name == "-1":
					continue
				var mesh_rotation = prototype["mesh_rotation"]
				spawn_mesh(mesh_name, Vector3(x, y, z), mesh_rotation)
