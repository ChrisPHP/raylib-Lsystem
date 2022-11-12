package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const width int32 = 1920
const height int32 = 1080

var (
	running = true
	output  = ""
)

func init() {
	rl.InitWindow(width, height, "L-system raylib")
	rl.SetTargetFPS(60)
	output = create_lsystem(6, "X")
	fmt.Println(output)
}

func update() {
	running = !rl.WindowShouldClose()
}

func render() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)
	draw_lsystem(output, 20)
	rl.DrawFPS(10, 10)
	rl.EndDrawing()
}

func main() {
	for running {
		update()
		render()
	}
	quit()
}

func quit() {
	rl.CloseWindow()
}

func create_lsystem(iterations int, axiom string) string {
	startString := axiom
	endString := ""
	for i := 0; i < iterations; i++ {
		endString = process_string(startString)
		startString = endString
	}

	return endString
}

func process_string(oldString string) string {
	newstr := ""
	for _, char := range oldString {
		newstr = newstr + apply_rules(string(char))
	}
	return newstr
}

func apply_rules(char string) string {
	newstr := ""
	if char == "F" {
		newstr = "FF"
	} else if char == "X" {
		newstr = "F[+X]F[-X]+X"
	} else {
		newstr = char
	}
	return newstr
}

func draw_lsystem(instructions string, angle float32) {
	PosX := float32(width / 2)
	PosY := float32(height)
	distance := float32(5)
	vertical_angle := float32(90)
	Stack := [][]float32{}
	for _, char := range instructions {
		cmd := string(char)
		if cmd == "F" {
			EndPosX, EndPosY := polar_to_cartesian(float64(distance), float64(vertical_angle))
			EndX := PosX + EndPosX
			EndY := PosY - EndPosY
			rl.DrawLineEx(rl.Vector2{X: float32(PosX), Y: float32(PosY)}, rl.Vector2{X: float32(EndX), Y: float32(EndY)}, 5, rl.Red)
			PosX = EndX
			PosY = EndY
		} else if cmd == "+" {
			vertical_angle += angle
		} else if cmd == "-" {
			vertical_angle -= angle
		} else if cmd == "[" {
			Stack = append(Stack, []float32{PosX, PosY, vertical_angle})
		} else if cmd == "]" {
			Popped := Stack[len(Stack)-1]
			PosX = Popped[0]
			PosY = Popped[1]
			vertical_angle = Popped[2]
			Stack = Stack[:len(Stack)-1]
		}
	}
}

func polar_to_cartesian(radians float64, vertical_angle float64) (float32, float32) {
	theta_vertical := vertical_angle * math.Pi / 180

	X := radians * math.Cos(theta_vertical)
	Y := radians * math.Sin(theta_vertical)

	return float32(X), float32(Y)
}
