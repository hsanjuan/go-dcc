package dcc

import (
	"testing"
)

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
