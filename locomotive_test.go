package dcc

import (
	"testing"

	"github.com/hsanjuan/go-dcc/driver/dummy"
)

func TestApply(t *testing.T) {
	l := &Locomotive{
		Name:        "loco",
		Address:     3,
		speedPacket: NewBroadcastIdlePacket(&dummy.DCCDummy{}),
	}
	l.Apply()
	if l.speedPacket != nil {
		t.Error("should have cleared the speed packet")
	}
}

func TestString(t *testing.T) {
	l := &Locomotive{
		Name:      "loco",
		Address:   4,
		Direction: Forward,
		Speed:     4,
		Fl:        true,
		F1:        true,
		F2:        true,
		F3:        true,
		F4:        true,
	}
	l.String()
}
