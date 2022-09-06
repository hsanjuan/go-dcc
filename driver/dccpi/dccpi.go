// package pi provides a Raspberry Pi driver for go-dcc.
//
// Note that the Raspberry Pi needs to be equipped with an additional booster
// circuit in order to send the signal to the tracks. See the README.md for
// more information.

package dccpi

import (
	"fmt"
	"os"

	rpio "github.com/stianeikeland/go-rpio/v4"
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

type DCCPi struct {
}

func NewDCCPi() (*DCCPi, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}
	return &DCCPi{}, nil
}

func (pi *DCCPi) Low() {
	SignalGPIO.Low()
}

func (pi *DCCPi) High() {
	SignalGPIO.High()
}

func (pi *DCCPi) TracksOff() {
	BrakeGPIO.High()
}

func (pi *DCCPi) TracksOn() {
	BrakeGPIO.Low()
}
