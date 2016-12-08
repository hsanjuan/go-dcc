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
	if c.GetLoco("abc") != nil {
		t.Error("loco should have been deleted")
	}
}

func TestStart(t *testing.T) {
	c := NewController(&dummy.DCCDummy{})
	c.Start()
	time.Sleep(1 * time.Second)
	c.AddLoco(&Locomotive{Name: "abc", Address: 10})
	time.Sleep(1 * time.Second)
	c.Stop()
}
