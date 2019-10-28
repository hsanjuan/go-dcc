// +build sam,atsamd21,arduino_nano33

package tinygo

import "machine"

// Default pins for Arduino_nano33 board
var (
	SignalPin = machine.D2
	BrakePin  = machine.D3
)
