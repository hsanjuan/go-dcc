package dccpi

import "testing"

func TestNew(t *testing.T) {
	d, err := NewDCCPi()
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
