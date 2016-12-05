package dcc

import "sync"

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

func (l *Locomotive) Apply() {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.speedPacket = nil
	l.flPacket = nil
}
