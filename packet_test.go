package dcc

import (
	"testing"
	"time"

	"github.com/hsanjuan/go-dcc/driver/dummy"
)

func dummyTesting() {
	dummy.ByteOneMax = 200 * time.Microsecond
	dummy.ByteZeroMax = 400 * time.Microsecond
	BitZeroPartDuration = 250 * time.Microsecond
	BitOnePartDuration = 150 * time.Microsecond
}

func TestSend(t *testing.T) {
	//dummyTesting()
	d := &dummy.DCCDummy{}
	p := NewBroadcastIdlePacket(d)
	d.TracksOn()
	p.Send()
	time.Sleep(1 * time.Second)
	packetStr := dummy.GuessBuffer.String()
	t.Log("Pckt: ", p.String())
	t.Log("Sent: ", packetStr)
	if packetStr != p.String() {
		t.Error("should have sent the encoded package")
	}
}

func TestNewPacket(t *testing.T) {
	p := NewPacket(&dummy.DCCDummy{}, 0xFF, []byte{0x01})
	if p.String() != "11111111111111110111111110000000010111111101" {
		t.Error("Bad packet: ", p.String())
	}
}

func TestNewBaselinePacket(t *testing.T) {
	p := NewBaselinePacket(&dummy.DCCDummy{}, 0xFF, []byte{0x01})
	if p.String() != "11111111111111110011111110000000010011111101" {
		t.Error("Bad packet: ", p.String())
	}
}

func TestIdlePacket(t *testing.T) {
	p := NewBroadcastIdlePacket(&dummy.DCCDummy{})
	if p.String() != "11111111111111110111111110000000000111111111" {
		t.Error("Bad idle packet")
	}
}

func TestNewSpeedAndDirectionPacket(t *testing.T) {
	p := NewSpeedAndDirectionPacket(&dummy.DCCDummy{}, 0xFF, 0xFF, Forward)
	if p.String() != "11111111111111110011111110011111110000000001" {
		t.Error("Bad speed and direction packet: ", p.String())
	}
}

func TestNewFunctionGroupOnePacket(t *testing.T) {
	p := NewFunctionGroupOnePacket(&dummy.DCCDummy{}, 0xFF, true, false, true, false, false)
	if p.String() != "11111111111111110111111110100100100011011011" {
		t.Error("Bad Function Group One packet: ", p.String())
	}
}

func TestNewBroadcastResetPacket(t *testing.T) {
	p := NewBroadcastResetPacket(&dummy.DCCDummy{})
	if p.String() != "11111111111111110000000000000000000000000001" {
		t.Error("Bad reset packet")
	}
}

func TestNewBroadcastStopPacket(t *testing.T) {
	p := NewBroadcastStopPacket(&dummy.DCCDummy{}, Backward, true, false)
	if p.String() != "11111111111111110000000000010000000010000001" {
		t.Error("Bad stop packet: ", p.String())
	}
}
