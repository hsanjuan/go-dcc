// package pi provides a Raspberry Pi driver for go-dcc.
//
// Note that the Raspberry Pi needs to be equipped with an additional booster
// circuit in order to send the signal to the tracks. See the README.md for
// more information.

package dccpi

import (
	"log"

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
		log.Fatal("Cannot initialize GPIO: ", err)
	}
	BrakeGPIO.Output()
	BrakeGPIO.Pull(rpio.PullUp)
	SignalGPIO.Output()

}

type DCCPi struct {
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
