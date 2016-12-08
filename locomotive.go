package dcc

import "sync"

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

	mux sync.Mutex

	speedPacket *DCCPacket
	flPacket    *DCCPacket
}

func (l *Locomotive) sendPackets(d DCCDriver) {
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
