// Package tinygo provides a TinyGo driver for use with go-dcc.
//
// The default pins used are defined in platform-specific files as
// variables named SignalPin and BrakePin. Not all TinyGo platforms have been
// added. In order to add an extra platform choose an existing platform file
// and adapt the build tags and the pins to those specific for that platform.
//
// In order to build these files you need TinyGo. For an example go-dcc
// controller using this driver see github.com/hsanjuan/go-dcc/examples/tinygo.
package tinygo

import (
	"machine"
)

// Driver implements a go-dcc driver for TinyGo
type Driver struct{}

// New returns a Driver after configuring the pins.
func New() Driver {
	SignalPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	BrakePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	BrakePin.High()
	return Driver{}
}

// Low sets the SignalPin to low.
func (tgo Driver) Low() {
	SignalPin.Low()
}

// High sets the SignalPin to high.
func (tgo Driver) High() {
	SignalPin.High()
}

// TracksOn sets the BrakePin to low.
func (tgo Driver) TracksOn() {
	BrakePin.Low()
}

// TracksOff sets the BrakePin to high.
func (tgo Driver) TracksOff() {
	BrakePin.High()
}
