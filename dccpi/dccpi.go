package main

import (
	"log"
	"time"

	dcc "github.com/hsanjuan/go-dcc"
	"github.com/hsanjuan/go-dcc/driver/dummy"
)

func main() {
	cfg, err := dcc.LoadConfig("testconfig.json")
	if err != nil {
		log.Fatal(err)
	}
	controller := dcc.NewControllerWithConfig(&dummy.DCCDummy{}, cfg)
	controller.Start()
	time.Sleep(20000 * time.Microsecond)
	controller.Stop()
}
