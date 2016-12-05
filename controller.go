package dcc

import (
	"fmt"
	"log"
)

// CommandRepeat specifies how many times a single
// packet is sent.
var CommandRepeat = 5

// CommandMaxQueue specifies how many commands can
// queue before sending a new command blocks
// the sender
var CommandMaxQueue = 3

type Controller struct {
	locomotives map[string]*Locomotive
	driver      DCCDriver

	doneCh     chan bool
	shutdownCh chan bool
	commandCh  chan *DCCPacket
}

func NewController(d DCCDriver) *Controller {
	return &Controller{
		driver:      d,
		locomotives: make(map[string]*Locomotive),
		doneCh:      make(chan bool),
		shutdownCh:  make(chan bool),
		commandCh:   make(chan *DCCPacket, CommandMaxQueue),
	}
}

func NewControllerWithConfig(d DCCDriver, cfg *DCCConfig) *Controller {
	c := NewController(d)

	for _, loco := range cfg.Locomotives {
		c.AddLoco(&Locomotive{
			Name:      loco.Name,
			Address:   loco.Address,
			Speed:     loco.Speed,
			Direction: loco.Direction,
			Fl:        loco.Fl})
	}
	return c
}

func (c *Controller) AddLoco(l *Locomotive) {
	c.locomotives[l.Name] = l
}

func (c *Controller) RmLoco(l *Locomotive) {
	delete(c.locomotives, l.Name)
}

func (c *Controller) GetLoco(n string) *Locomotive {
	return c.locomotives[n]
}

func (c *Controller) Start() {
	go c.run()
}

func (c *Controller) Stop() {
	log.Println("Stopping controller")
	close(c.shutdownCh)
	fmt.Println("ch sent")
	<-c.doneCh
	log.Println("All Stop. Tracks Off.")
}

func (c *Controller) run() {
	idle := NewBroadcastIdlePacket(c.driver)
	stop := NewBroadcastStopPacket(c.driver, Forward, false, true)
	c.driver.TracksOn()
	for {
		select {
		case <-c.shutdownCh:
			stop.Send()
			c.driver.TracksOff()
			defer close(c.doneCh)
			return
		case p := <-c.commandCh:
			for i := 0; i < CommandRepeat; i++ {
				p.Send()
			}
		default:
			if len(c.locomotives) == 0 {
				c.commandCh <- idle
				break
			}

			for _, loco := range c.locomotives {
				loco.sendPackets(c.driver)
			}
			idle.PacketSpace()
		}
	}
}
