# 8BitDo Zero 2 Wrapper

This wrapper takes the input from the joystick to turn it into an effective remote controller for my Hyprland setup. It is totally based off my setup, so it probably won't work in your computer.

## Roadmap

**Only working in Linux for the foreseeable future**

- [x] Detect input
- [x] Map inputs to OS commands
- [ ] Handle joystick connection
- [ ] Handle joystick disconnection
- [ ] Handle multiple input layers
- [ ] Parse `.json` config

## DevLog

### 2024-05-22

Still using `bluetoothctl` and `btmon` to manually connect the joystick and debug it. The next step is definitely handling the joystick connection and disconnection automatically (I will probably look at [TinyGo Bluetooth Library](https://github.com/tinygo-org/bluetooth) to get a feel for how they handle this).

