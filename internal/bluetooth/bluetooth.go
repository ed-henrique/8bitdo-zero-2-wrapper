package bluetooth

import (
	"errors"
	"fmt"
	"os/exec"
)

const (
	newLine           = 10
	carriageReturn    = 13
	defaultDeviceName = "8BitDo Zero 2 gamepad"
)

func Scan(timeout uint64) error {
	cmd := exec.Command("bluetoothctl", "--timeout", fmt.Sprintf("%d", timeout), "scan", "on")
	return cmd.Run()
}

func JoystickDeviceAddress() (string, error) {
	cmd := exec.Command("bluetoothctl", "devices")
	devicesRaw, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	var newLineIndex int
	for i, d := range devicesRaw {
		if d == newLine || d == carriageReturn {
			if string(devicesRaw[newLineIndex+25:i]) == defaultDeviceName {
				return string(devicesRaw[newLineIndex+7 : newLineIndex+24]), nil
			}
			newLineIndex = i + 1
		}
	}

	return "", errors.New("joystick not found")
}

func ConnectJoystickDevice(address string) error {
	cmd := exec.Command("bluetoothctl", "pair", address)
	err := cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("bluetoothctl", "connect", address)
	err = cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("bluetoothctl", "trust", address)
	return cmd.Run()
}

func IsJoystickReady(address string) error {
	cmd := exec.Command("bluetoothctl", "info", address)
	infoRaw, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}

	var (
		newLineCount    int
		newLineOldIndex int
		newLineNewIndex int
	)
	for i, r := range infoRaw {
		if r == newLine || r == carriageReturn {
			newLineCount++
			newLineOldIndex = newLineNewIndex
			newLineNewIndex = i

			switch newLineCount {
			case 6:
				if string(infoRaw[newLineOldIndex+10:newLineNewIndex]) != "yes" {
					return errors.New("joystick not paired")
				}
			case 8:
				if string(infoRaw[newLineOldIndex+11:newLineNewIndex]) != "yes" {
					return errors.New("joystick not trusted")
				}
			case 10:
				if string(infoRaw[newLineOldIndex+13:newLineNewIndex]) != "yes" {
					return errors.New("joystick not connected")
				}
			case 11:
				return nil
			}
		}
	}

	return nil
}

func DisconnectJoystickDevice(address string) error {
	cmd := exec.Command("bluetoothctl", "remove", address)
	return cmd.Run()
}
