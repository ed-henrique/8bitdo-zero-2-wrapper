package joystick

import (
	"8bitdo-zero-2-wrapper/internal/errkit"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
)

// Last 4 bits tell a lot about the button being pressed:
//
// - First bye is related to the button being pressed or released
// - Second byte is related to directional buttons, but I can't get a feel for what they actually do
// - Third byte is related to the type of button being pressed (directional or not)
// - Fourth byte is related to:
//   - The button ID if not directional
//   - The axis if directional

const (
	configDefaultPath = "./config.json"
	defaultDevicePath = "/dev/input/js0"

	buttonPressedIndex = 4
	isButtonReleased   = 0

	typeIndex       = 6
	typeDirectional = 2

	directionIndex = 4
	axisIndex      = 7
	axisXIndex     = 0
	axisYIndex     = 1

	buttonIDIndex = 7
	idA           = 0
	idB           = 1
	idX           = 3
	idY           = 4
	idR           = 7
	idL           = 6
	idStart       = 11
	idSelect      = 10
	idLeft        = 1
	idRight       = 255
	idUp          = 1
	idDown        = 255

	noButton = iota
	buttonA
	buttonB
	buttonX
	buttonY
	buttonR
	buttonL
	buttonStart
	buttonSelect
	buttonLeft
	buttonRight
	buttonUp
	buttonDown
)

type button uint8

func (b button) String() string {
	switch b {
	case buttonA:
		return "A"
	case buttonB:
		return "B"
	case buttonX:
		return "X"
	case buttonY:
		return "Y"
	case buttonR:
		return "R"
	case buttonL:
		return "L"
	case buttonStart:
		return "START"
	case buttonSelect:
		return "SELECT"
	case buttonLeft:
		return "LEFT"
	case buttonRight:
		return "RIGHT"
	case buttonUp:
		return "UP"
	case buttonDown:
		return "DOWN"
	default:
		return ""
	}
}

type actionConfig map[string]string

func newActionConfigJSON(path string) actionConfig {
	configData, err := os.ReadFile(path)
	errkit.Must(err)

	var c actionConfig
	err = json.Unmarshal(configData, &c)
	errkit.Must(err)

	return c
}

type joystickConfig struct {
	configPath string
	devicePath string
}

type joystickOption interface {
	apply(*joystickConfig)
}

type joystickConfigOption func(*joystickConfig)

func (o joystickConfigOption) apply(c *joystickConfig) {
	o(c)
}

func WithConfigPath(configPath string) joystickConfigOption {
	return joystickConfigOption(func(o *joystickConfig) {
		o.configPath = configPath
	})
}

func WithDevicePath(devicePath string) joystickConfigOption {
	return joystickConfigOption(func(o *joystickConfig) {
		o.devicePath = devicePath
	})
}

type Joystick struct {
	actionConfigPath string
	actionConfig     actionConfig
	devicePath       string
	deviceFile       *os.File
	raw              []byte
	button           button
	pressed          bool
}

func New(opts ...joystickConfigOption) Joystick {
	c := joystickConfig{
		devicePath: defaultDevicePath,
		configPath: configDefaultPath,
	}

	for _, opt := range opts {
		opt.apply(&c)
	}

	cfg := newActionConfigJSON(c.configPath)

	return Joystick{
		actionConfig: cfg,
		devicePath:   c.devicePath,
		deviceFile:   nil,
		raw:          make([]byte, 8),
		button:       noButton,
		pressed:      false,
	}
}

func (j *Joystick) Start() {
	if j.deviceFile != nil {
		return
	}

	deviceRawData, err := os.Open(j.devicePath)
	errkit.Must(err)

	j.deviceFile = deviceRawData
}

func (j *Joystick) Stop() {
	j.deviceFile.Close()
}

func (j Joystick) runButtonCommand() error {
	if !j.pressed {
		return nil
	}

	v, ok := j.actionConfig[j.button.String()]

	if !ok {
		return errors.New("Button being pressed was not detected in map.")
	}

	// TODO: Make this arg parsing better
	args := strings.Split(v, " ")
	cmd := exec.Command(args[0], args[1:]...)
	return cmd.Run()
}

func (j *Joystick) HandleInput() error {
	n, err := j.deviceFile.Read(j.raw)

	if err != nil {
		return err
	}

	if n > 0 {
		j.pressed = j.raw[buttonPressedIndex] != isButtonReleased

		switch j.raw[typeIndex] {
		case typeDirectional:
			switch j.raw[axisIndex] {
			case axisXIndex:
				switch j.raw[directionIndex] {
				case idLeft:
					j.button = buttonLeft
				case idRight:
					j.button = buttonRight
				}
			case axisYIndex:
				switch j.raw[directionIndex] {
				case idUp:
					j.button = buttonUp
				case idDown:
					j.button = buttonDown
				}
			}
		default:
			switch j.raw[buttonIDIndex] {
			case idA:
				j.button = buttonA
			case idB:
				j.button = buttonB
			case idX:
				j.button = buttonX
			case idY:
				j.button = buttonY
			case idR:
				j.button = buttonR
			case idL:
				j.button = buttonL
			case idStart:
				j.button = buttonStart
			case idSelect:
				j.button = buttonSelect
			}
		}

		err = j.runButtonCommand()

		if err != nil {
			return err
		}
	}

	return nil
}
