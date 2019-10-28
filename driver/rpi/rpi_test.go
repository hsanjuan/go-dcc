package rpi

import "testing"

func TestNew(t *testing.T) {
	d, err := NewRPi()
	if d == nil && err == nil {
		t.Error("cannot return both nil")
	}

	if err == nil {
		d.TracksOn()
		d.Low()
		d.High()
		d.TracksOff()
	}
}
