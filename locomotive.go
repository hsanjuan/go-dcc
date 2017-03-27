package dcc

import (
	"fmt"
	"sync"
)

// Direction constants.
const (
	Backward Direction = 0
	Forward  Direction = 1
)

// Direction represents the locomotive direction and can be
// Forward or Backward.
type Direction byte

// Locomotive represents a DCC device, usually a locomotive.
// Locomotives are represented by their name and address and
// include certain properties like speed, direction or FL.
// Each locomotive produces two packets: one speed and direction
// packet and one Function Group One packet.
// Block can be used if the user wishes to integrate with a
// custom signalling system, to track the location of the loco
type Locomotive struct {
	Name      string    `json:"name"`
	Address   uint8     `json:"address"`
	Speed     uint8     `json:"speed"`
	Direction Direction `json:"direction"`
	Fl        bool      `json:"fl"`
	F1        bool      `json:"f1"`
	F2        bool      `json:"f2"`
	F3        bool      `json:"f3"`
	F4        bool      `json:"f4"`
	Block			string		`json:"block"`

	mux sync.Mutex

	speedPacket *Packet
	flPacket    *Packet
}

func (l *Locomotive) String() string {
	var dir, fl, f1, f2, f3, f4 string = "", "off", "off", "off", "off", "off"
	if l.Direction == Forward {
		dir = ">"
	} else {
		dir = "<"
	}
	if l.Fl {
		fl = "on"
	}
	if l.F1 {
		f1 = "on"
	}
	if l.F2 {
		f2 = "on"
	}
	if l.F3 {
		f3 = "on"
	}
	if l.F4 {
		f4 = "on"
	}
	return fmt.Sprintf("%s:%d |%d%s| |%s| |%s|%s|%s|%s|",
		l.Name,
		l.Address,
		l.Speed,
		dir,
		fl,
		f1,
		f2,
		f3,
		f4)
}

func (l *Locomotive) sendPackets(d Driver) {
	if l.speedPacket == nil {
		l.mux.Lock()
		l.speedPacket = NewSpeedAndDirectionPacket(d,
			l.Address, l.Speed, l.Direction)
		l.mux.Unlock()
	}
	if l.flPacket == nil {
		l.mux.Lock()
		l.flPacket = NewFunctionGroupOnePacket(d,
			l.Address, l.Fl, l.F1, l.F2, l.F3, l.F4)
		l.mux.Unlock()
	}
	l.speedPacket.Send()
	l.flPacket.Send()
}

// Apply makes any changes to the Locomotive's properties
// to be reflected in the packets generated for it and,
// therefore, alter the behaviour of the device on the tracks.
func (l *Locomotive) Apply() {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.speedPacket = nil
	l.flPacket = nil
}
