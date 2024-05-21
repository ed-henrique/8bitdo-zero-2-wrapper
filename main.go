package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	JOYSTICK_DEVICE = "/dev/input/js0"

	TYPE        = 6
	DIRECTIONAL = 2

	PRESS    = 4
	RELEASED = 0

	AXIS      = 7
	X_AXIS    = 0
	Y_AXIS    = 1
	DIRECTION = 4

	BUTTON = 7

	RAW_A      = 0
	RAW_B      = 1
	RAW_X      = 3
	RAW_Y      = 4
	RAW_R      = 7
	RAW_L      = 6
	RAW_START  = 11
	RAW_SELECT = 10
	RAW_LEFT   = 1
	RAW_RIGHT  = 255
	RAW_UP     = 1
	RAW_DOWN   = 255

	A = iota + 1
	B
	X
	Y
	R
	L
	START
	SELECT
	LEFT
	RIGHT
	UP
	DOWN
)

type Button uint8

func (b Button) String() string {
	switch b {
	case A:
		return "A"
	case B:
		return "B"
	case X:
		return "X"
	case Y:
		return "Y"
	case R:
		return "R"
	case L:
		return "L"
	case START:
		return "START"
	case SELECT:
		return "SELECT"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	default:
		return ""
	}
}

type Action uint8

var ButtonAction map[Button]Action

type Joystick struct {
	Raw     []byte
	Button  Button
	Pressed bool
}

func (j Joystick) Run() error {
	if !j.Pressed {
		return nil
	}

	var cmd *exec.Cmd
	switch j.Button {
	case A:
		cmd = exec.Command("")
	case B:
		cmd = exec.Command("")
	case X:
		cmd = exec.Command("")
	case Y:
		cmd = exec.Command("")
	case R:
		cmd = exec.Command("")
	case L:
		cmd = exec.Command("")
	case START:
		cmd = exec.Command("")
	case SELECT:
		cmd = exec.Command("")
	case LEFT:
		cmd = exec.Command("")
	case RIGHT:
		cmd = exec.Command("")
	case UP:
		cmd = exec.Command("")
	case DOWN:
		cmd = exec.Command("")
	default:
		return nil
	}

	return cmd.Run()
}

// Last 4 bits tell a lot about the button being pressed:
//
// - First bye is related to the button being pressed or released
// - Second byte is related to directional buttons, but I can't get a feel for what they actually do
// - Third byte is related to the type of button being pressed (directional or not)
// - Fourth byte is related to:
//   - The button ID if not directional
//   - The axis if directional

func main() {
	joyRaw, err := os.Open(JOYSTICK_DEVICE)

	if err != nil {
		panic(err)
	}

	defer joyRaw.Close()
	joystick := Joystick{
		Raw: make([]byte, 8),
	}

	for {
		n, err := joyRaw.Read(joystick.Raw)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if n > 0 {
			joystick.Pressed = joystick.Raw[PRESS] != RELEASED

			switch joystick.Raw[TYPE] {
			case DIRECTIONAL:
				switch joystick.Raw[AXIS] {
				case X_AXIS:
					switch joystick.Raw[DIRECTION] {
					case RAW_LEFT:
						joystick.Button = LEFT
					case RAW_RIGHT:
						joystick.Button = RIGHT
					}
				case Y_AXIS:
					switch joystick.Raw[DIRECTION] {
					case RAW_UP:
						joystick.Button = UP
					case RAW_DOWN:
						joystick.Button = DOWN
					}
				}
			default:
				switch joystick.Raw[BUTTON] {
				case RAW_A:
					joystick.Button = A
				case RAW_B:
					joystick.Button = B
				case RAW_X:
					joystick.Button = X
				case RAW_Y:
					joystick.Button = Y
				case RAW_R:
					joystick.Button = R
				case RAW_L:
					joystick.Button = L
				case RAW_START:
					joystick.Button = START
				case RAW_SELECT:
					joystick.Button = SELECT
				}
			}

			fmt.Printf("Button: %s Pressed: %t\n", joystick.Button, joystick.Pressed)
		}
	}
}
