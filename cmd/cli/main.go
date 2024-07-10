package main

import (
	"8bitdo-zero-2-wrapper/internal/bluetooth"
	"8bitdo-zero-2-wrapper/internal/errkit"
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls", "/dev/input")
	devicesRaw, errLs := cmd.CombinedOutput()
	if errLs != nil {
		fmt.Println(errLs)
	}

	fmt.Println(string(devicesRaw))

	address, err := bluetooth.JoystickDeviceAddress()
	if err == nil {
		err = bluetooth.IsJoystickReady(address)
		if err != nil {
			fmt.Println(err)
			err = bluetooth.ConnectJoystickDevice(address)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println(err)
	}

	for ;err != nil; err = bluetooth.IsJoystickReady(address) {
		err = bluetooth.Scan(5)
		errkit.Must(err)

		address, err = bluetooth.JoystickDeviceAddress()
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = bluetooth.ConnectJoystickDevice(address)
		if err != nil {
			fmt.Println(err)
		}
	}

	cmd2 := exec.Command("ls", "/dev/input")
	devicesRaw2, errLs2 := cmd2.CombinedOutput()
	if errLs2 != nil {
		fmt.Println(errLs2)
	}

	fmt.Println(string(devicesRaw2))

	fmt.Printf("Funcionou e o endereço é %s\n", address)

	/*
	err = bluetooth.Scan(30)
	if err != nil {
		fmt.Println(err)
	}

	address, err := bluetooth.JoystickDeviceAddress()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(address)

	err = bluetooth.ConnectJoystickDevice(address)
	if err != nil {
		fmt.Println(err)
	}

	/*
	j := joystick.New()
	j.Start()
	defer j.Stop()

	var err error
	for {
		err = j.HandleInput()

		if err != nil {
			fmt.Println(err)
		}
	}
	*/
}
