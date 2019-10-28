package main

// This example shows how to create an embeddable DCC controller using TinyGo
// to build (and flash) a platform-specific build that can be used on many
// microcontrollers.
//
// You may need to add your platforms or adapt the pins used. Not all
// TinyGo-supported platforms have been added to the DCC driver. See the
// documentation for github.com/hsanjuan/go-dcc/tinygo/driver/tinygo for more
// information.
import (
	_ "machine"

	"github.com/hsanjuan/go-dcc"
	"github.com/hsanjuan/go-dcc/driver/tinygo"
)

func main() {
	// Customize pins - warning, platform dependent!
	// tinygo.SignalPin = machine.DP3
	// tinygo.BrakePin = machine.DP4

	d := tinygo.Driver{}
	ctrl := dcc.NewController(d)
	ctrl.AddLoco(dcc.Locomotive{
		Name:    "addr1",
		Address: 5,
		Speed:   10,
	})
	d.TracksOn()
	for {
		// Controller logic

		c.SendPackets()
	}
	d.TracksOff()
}
