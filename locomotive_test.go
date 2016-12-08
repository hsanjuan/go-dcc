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
