package main

import (
	"log"
	"time"

	dcc "github.com/hsanjuan/go-dcc"
	"github.com/hsanjuan/go-dcc/driver/dccpi"
)

func main() {
	cfg, err := dcc.LoadConfig("testconfig.json")
	if err != nil {
		log.Fatal(err)
	}
	controller := dcc.NewControllerWithConfig(&dccpi.DCCPi{}, cfg)
	//controller := dcc.NewControllerWithConfig(&dummy.DCCDummy{}, cfg)
	controller.Start()
	time.Sleep(60 * time.Second)
	controller.Stop()
}
