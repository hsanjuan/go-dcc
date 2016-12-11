package dcc

import (
	"testing"
	"time"

	"github.com/hsanjuan/go-dcc/driver/dummy"
)

func TestNewController(t *testing.T) {
	cfg, _ := LoadConfig("./test/config.json")
	c := NewControllerWithConfig(&dummy.DCCDummy{}, cfg)
	c.Stop()
}

func TestAddLoco(t *testing.T) {
	c := NewController(&dummy.DCCDummy{})
	c.AddLoco(&Locomotive{})
}

func TestRmLoco(t *testing.T) {
	c := NewController(&dummy.DCCDummy{})
	c.AddLoco(&Locomotive{Name: "abc"})
	c.RmLoco(&Locomotive{Name: "abc"})
	_, ok := c.GetLoco("abc")
	if ok {
		t.Error("loco should have been deleted")
	}
}

func TestLocos(t *testing.T) {
	c := NewController(&dummy.DCCDummy{})
	c.AddLoco(&Locomotive{Name: "abc"})
	l := c.Locos()
	if len(l) != 1 || l[0].Name != "abc" {
		t.Error("Locos() does not work")
	}
}

func TestCommand(t *testing.T) {
	d := &dummy.DCCDummy{}
	c := NewController(d)
	p := NewBroadcastIdlePacket(d)
	c.Command(p)
	c.Start()
	time.Sleep(250 * time.Millisecond)
	c.Stop()
}

func TestStart(t *testing.T) {
	c := NewController(&dummy.DCCDummy{})
	c.Start()
	time.Sleep(1 * time.Second)
	c.AddLoco(&Locomotive{Name: "abc", Address: 10})
	time.Sleep(1 * time.Second)
	c.Stop()
}
