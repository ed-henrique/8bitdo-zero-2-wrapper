package main

import (
	"8bitdo-zero-2-wrapper/internal/joystick"
	"fmt"
)

func main() {
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
}
