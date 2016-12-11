package dummy

import (
	"testing"
	"time"
)

func TestGuessBuffer(t *testing.T) {
	d := DCCDummy{}
	d.TracksOn()
	d.Low()
	d.High()
	d.Low()
	d.High()
	if GuessBuffer.String() != "11" {
		t.Error("it should guess 1")
	}
	d.TracksOff()
	d.TracksOn()
	d.Low()
	time.Sleep(5000 * time.Microsecond)
	d.High()
	time.Sleep(5000 * time.Microsecond)
	d.Low()
	time.Sleep(5000 * time.Microsecond)
	d.High()
	d.Low()
	time.Sleep(time.Second)
	d.High()

	if GuessBuffer.String() != "00\n" {
		t.Error("it should guess 0")
	}
}
