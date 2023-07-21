extends Node

@onready var module_scene = preload("res://Module.tscn").instantiate()
@onready var modules_scene = preload("res://wfc_modules.glb").instantiate()

const MAP_FILE_NAME = "map.json"

# Assuming this node has a function spawn_mesh
func spawn_mesh(mesh_name: String, position: Vector3, rotation: int):
	#print(mesh_name)
	
	var inst = module_scene.duplicate()
	add_child(inst)
	# meshes.append(inst) # TODO - keep track and free
	inst.transform = Transform3D(Basis(Vector3(0, 1, 0), rotation), position)

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
	var prototypes = json_dict["Prototypes"]
	for prototype in prototypes:
		var mesh_name = prototype["mesh_name"]
		if mesh_name == "-1":
			continue
		var mesh_rotation = prototype["mesh_rotation"]
		var position = Vector3(prototype["position"]["X"], prototype["position"]["Y"], prototype["position"]["Z"])
		spawn_mesh(mesh_name, position, mesh_rotation)
