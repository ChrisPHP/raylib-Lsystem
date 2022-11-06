package main

import (
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const width int32 = 1920
const height int32 = 1080

var (
	running    = true
	state      = 1
	output     = ""
	DrawButton = rl.Rectangle{
		X:      1700,
		Y:      50,
		Width:  200,
		Height: 50,
	}
	ClearButton = rl.Rectangle{
		X:      1700,
		Y:      100,
		Width:  200,
		Height: 50,
	}
	camera = rl.Camera3D{
		Position: rl.Vector3{
			X: 10,
			Y: 10,
			Z: 10,
		},
		Target: rl.Vector3{
			X: 0,
			Y: 0,
			Z: 0,
		},
		Up: rl.Vector3{
			X: 0,
			Y: 1,
			Z: 0,
		},
		Fovy:       90,
		Projection: rl.CameraPerspective,
	}
)

func drawScene() {
	rl.DrawText(output, 0, 1000, 20, rl.Black)
	draw_lsystem(output, 90, 0.5)
}

func input(mouse rl.Vector2) {
	if rl.CheckCollisionPointRec(mouse, DrawButton) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			state = 1
			output = create_lsystem(5, "X")
		}
	}
	if rl.CheckCollisionPointRec(mouse, ClearButton) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			output = ""
			state = 2
		}
	}
	if rl.IsKeyPressed(rl.KeyF) {
		state = 1
		output = create_lsystem(5, "X")
	} else if rl.IsKeyPressed(rl.KeyR) {
		output = ""
		state = 2
	}
}
func update() {
	running = !rl.WindowShouldClose()
}
func render() {
	rl.BeginDrawing()
	rl.BeginMode3D(camera)
	if state == 1 {
		drawScene()
	}
	rl.DrawGrid(10, 1)
	rl.ClearBackground(rl.RayWhite)
	rl.EndMode3D()

	rl.DrawRectangleRec(DrawButton, rl.DarkBlue)
	rl.DrawRectangleRec(ClearButton, rl.Lime)
	rl.DrawText("Draw Tree", 1750, 65, 20, rl.White)
	rl.DrawText("Clear Tree", 1750, 115, 20, rl.White)
	rl.DrawFPS(10, 10)
	rl.EndDrawing()
}

func init() {
	state = 1
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(width, height, "L-system raylib")
	rl.SetCameraMode(camera, rl.CameraFree)
	rl.SetTargetFPS(60)
}

func quit() {
	rl.CloseWindow()
}

func main() {
	mouse := rl.Vector2{X: 0, Y: 0}
	for running {
		rl.UpdateCamera(&camera)
		mouse = rl.GetMousePosition()
		input(mouse)
		update()
		render()
	}

	quit()
}

func apply_rules(char string) string {
	newstr := ""
	if char == "F" {
		newstr = "FF"
	} else if char == "X" {
		rand_num := rand.Intn(100-1+1) + 1
		if rand_num <= 30 {
			newstr = "F[+X]F[-X]+X"
		} else if rand_num > 30 && rand_num < 60 {
			newstr = "/F[+X]F[-X]+X"
		} else if rand_num >= 60 {
			newstr = "F[+X]dF[/-X]+X"
		}
	} else if char == "A" {
		newstr = "-F-FF-F"
	} else {
		newstr = char
	}
	return newstr
}

func process_string(oldString string) string {
	newstr := ""
	for _, char := range oldString {
		newstr = newstr + apply_rules(string(char))
	}
	return newstr
}

func polar_to_cartesian(radian float32, horizontal_angle float32, vertical_angle float32) (float32, float32, float32) {
	theta_hori := horizontal_angle * math.Pi / 180
	theta_verti := vertical_angle * math.Pi / 180
	X := float64(radian) * math.Cos(float64(theta_verti)) * math.Cos(float64(theta_hori))
	Y := float64(radian) * math.Sin(float64(theta_verti)) * math.Sin(float64(theta_hori))
	Z := float64(radian) * math.Cos(float64(theta_verti))
	return float32(X), float32(Y), float32(Z)
}

func draw_lsystem(instructions string, angle float32, distance float32) {
	PosX := float32(0)
	PosY := float32(0)
	PosZ := float32(0)
	vertical_angle := angle
	horizontal_angle := angle
	stack := [][]float32{}
	for _, char := range instructions {
		cmd := string(char)
		if cmd == "F" {
			EndPosX, EndPosY, EndPosZ := polar_to_cartesian(distance, horizontal_angle, vertical_angle)
			EndX := PosX + EndPosX
			EndY := PosY + EndPosY
			EndZ := PosZ + EndPosZ
			//rl.DrawLineEx(rl.Vector2{X: PosX, Y: PosY}, rl.Vector2{X: EndX,Y: EndY}, 5, rl.Red)
			rl.DrawCylinderEx(rl.Vector3{X: PosX, Y: PosY, Z: PosZ}, rl.Vector3{X: EndX, Y: EndY, Z: EndZ}, float32(0.2), float32(0.2), 4, rl.Red)
			//rl.DrawLine3D(rl.Vector3{X: PosX, Y: PosY, Z: PosZ}, rl.Vector3{X: EndX, Y: EndY, Z: EndZ}, rl.Red)
			PosX = EndX
			PosY = EndY
			PosZ = EndZ
		} else if cmd == "+" {
			vertical_angle += 20
		} else if cmd == "-" {
			vertical_angle -= 20
		} else if cmd == "/" {
			horizontal_angle += 20
		} else if cmd == "d" {
			horizontal_angle -= 20
		} else if cmd == "[" {
			stack = append(stack, []float32{PosX, PosY, PosZ, vertical_angle, horizontal_angle, distance})
		} else if cmd == "]" {
			popped := stack[len(stack)-1]
			PosX = popped[0]
			PosY = popped[1]
			PosZ = popped[2]
			vertical_angle = popped[3]
			horizontal_angle = popped[4]
			distance = popped[5]
			stack = stack[:len(stack)-1]
		}
	}
}

func create_lsystem(iterations int, axiom string) string {
	startString := axiom
	endstring := ""
	for i := 0; i < iterations; i++ {
		endstring = process_string(startString)
		startString = endstring
	}

	return endstring
}
