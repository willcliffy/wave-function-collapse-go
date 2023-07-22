extends Node3D

@onready var Camera = $Camera

const CAMERA_ROTATION_SPEED = 50.0

func _process(delta):
	rotate_y(delta / 4.0)

func _input(event):
	if not event is InputEventMouseButton: return
	handle_camera_zoom_input(event)


func _physics_process(delta):
	if Input.is_anything_pressed():
		handle_camera_rotation_input(delta)

func handle_camera_zoom_input(event):
	if event.button_index == MOUSE_BUTTON_WHEEL_UP and Camera.position.z > 5:
		Camera.position.z -= 1
		return
	if event.button_index == MOUSE_BUTTON_WHEEL_DOWN and Camera.position.z < 500:
		Camera.position.z += 1
		return

func handle_camera_rotation_input(delta):
	var rotate_horizontal := Input.get_action_strength("move_right") - Input.get_action_strength("move_left")
	var rotate_vertical := Input.get_action_strength("move_back") - Input.get_action_strength("move_forward")
	if rotate_horizontal != 0 or rotate_vertical != 0:
		var new_rotation = rotation_degrees
		new_rotation.y += rotate_horizontal * delta * CAMERA_ROTATION_SPEED
		new_rotation.x = new_rotation.x + rotate_vertical * delta * CAMERA_ROTATION_SPEED
		rotation_degrees = new_rotation
