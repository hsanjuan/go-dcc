package dcc

import (
	"testing"
	"time"

	"github.com/hsanjuan/go-dcc/driver/dummy"
)

var testCfg = Config{
	Locomotives: []Locomotive{
		{
			Name:      "Loco1",
			Address:   6,
			Speed:     5,
			Direction: 0,
			Fl:        true,
			F1:        false,
			F2:        false,
			F3:        false,
			F4:        false,
		},
		{
			Name:      "Loco2",
			Address:   5,
			Speed:     0,
			Direction: 0,
			Fl:        false,
			F1:        false,
			F2:        false,
			F3:        false,
			F4:        false,
		},
		{
			Name:      "loco3",
			Address:   4,
			Speed:     0,
			Direction: 0,
			Fl:        false,
			F1:        false,
			F2:        false,
			F3:        false,
			F4:        false,
		},
	},
}

func TestNewController(t *testing.T) {
	c := NewControllerWithConfig(&dummy.Driver{}, testCfg)
	c.Stop()
}

func TestAddLoco(t *testing.T) {
	c := NewController(&dummy.Driver{})
	c.AddLoco(Locomotive{})
}

func TestRmLoco(t *testing.T) {
	c := NewController(&dummy.Driver{})
	c.AddLoco(Locomotive{Name: "abc"})
	c.RmLoco("abc")
	_, ok := c.GetLoco("abc")
	if ok {
		t.Error("loco should have been deleted")
	}
}

func TestLocos(t *testing.T) {
	c := NewController(&dummy.Driver{})
	c.AddLoco(Locomotive{Name: "abc"})
	l := c.Locos()
	if len(l) != 1 || l[0].Name != "abc" {
		t.Error("Locos() does not work")
	}
}

func TestCommand(t *testing.T) {
	d := &dummy.Driver{}
	c := NewController(d)
	p := NewBroadcastIdlePacket(d)
	c.Command(p)
	c.Start()
	time.Sleep(250 * time.Millisecond)
	c.Stop()
}

func TestStart(t *testing.T) {
	c := NewController(&dummy.Driver{})
	c.Start()
	time.Sleep(1 * time.Second)
	c.AddLoco(Locomotive{Name: "abc", Address: 10})
	time.Sleep(1 * time.Second)
	c.Stop()
}
