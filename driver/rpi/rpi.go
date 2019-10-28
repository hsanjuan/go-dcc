// Package rpi provides a Raspberry Pi driver for go-dcc.
//
// Note that the Raspberry Pi needs to be equipped with an additional booster
// circuit in order to send the signal to the tracks. See the README.md for
// more information.
package rpi

import (
	"fmt"
	"os"

	rpio "github.com/stianeikeland/go-rpio"
)

// GPIO Outputs for the Raspberry PI DCC encoder
var (
	BrakeGPIO  rpio.Pin = 27
	SignalGPIO rpio.Pin = 17
)

func init() {
	err := rpio.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot initialize GPIO: "+err.Error()+".\n")
		return
	}
	BrakeGPIO.Output()
	BrakeGPIO.Pull(rpio.PullUp)
	SignalGPIO.Output()
}

// Driver implements a driver that can control Raspberry Pi GPIO pins.
type Driver struct {
}

// New returns a new driver. It will return an error if gpio is not
// accessible.
func New() (Driver, error) {
	err := rpio.Open()
	if err != nil {
		return Driver{}, err
	}
	return Driver{}, nil
}

// Low sets the SignalGPIO pin to low.
func (pi Driver) Low() {
	SignalGPIO.Low()
}

// High sets the SignalGPIO pin to high.
func (pi Driver) High() {
	SignalGPIO.High()
}

// TracksOff sets the BrakeGPIO pin to high.
func (pi Driver) TracksOff() {
	BrakeGPIO.High()
}

// TracksOn sets the BrakeGPIO pin to low.
func (pi Driver) TracksOn() {
	BrakeGPIO.Low()
}
