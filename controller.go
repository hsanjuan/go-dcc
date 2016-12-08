package dcc

import (
	"log"
	"sync"
)

// CommandRepeat specifies how many times a single
// packet is sent.
var CommandRepeat = 5

// CommandMaxQueue specifies how many commands can
// queue before sending a new command blocks
// the sender
var CommandMaxQueue = 3

// Controller represents a DCC Control Station. The
// controller keeps tracks of the DCC Locomotives and
// is in charge of sending DCC packets continuously to
// the tracks.
type Controller struct {
	locomotives map[string]Locomotive
	mux         sync.Mutex
	driver      Driver

	started    bool
	doneCh     chan bool
	shutdownCh chan bool
	commandCh  chan *Packet
}

// NewController builds a Controller.
func NewController(d Driver) *Controller {
	return &Controller{
		driver:      d,
		locomotives: make(map[string]Locomotive),
		doneCh:      make(chan bool),
		shutdownCh:  make(chan bool),
		commandCh:   make(chan *Packet, CommandMaxQueue),
	}
}

// NewControllerWithConfig builds a new Controller using the
// given configuration.
func NewControllerWithConfig(d Driver, cfg *Config) *Controller {
	c := NewController(d)

	for _, loco := range cfg.Locomotives {
		c.AddLoco(Locomotive{
			Name:      loco.Name,
			Address:   loco.Address,
			Speed:     loco.Speed,
			Direction: loco.Direction,
			Fl:        loco.Fl})
	}
	return c
}

// AddLoco adds a DCC device to the controller. The device
// will start receiving packets if the controller is running.
func (c *Controller) AddLoco(l Locomotive) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.locomotives[l.Name] = l
}

// RmLoco removes a DCC device from the controller. There
// will be no longer packets sent to it.
func (c *Controller) RmLoco(l Locomotive) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.locomotives, l.Name)
}

// GetLoco retrieves a DCC device by its Name. The boolean is
// true if the Locomotive was found.
func (c *Controller) GetLoco(n string) (Locomotive, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	l, ok := c.locomotives[n]
	return l, ok
}

// Start starts the controller: powers on the tracks
// and starts sending packets on them.
func (c *Controller) Start() {
	go c.run()
	c.started = true
}

// Stop shuts down the controller by stopping to send
// packets and removing power from the tracks.
func (c *Controller) Stop() {
	if c.started {
		c.shutdownCh <- true
		<-c.doneCh
		c.started = false
	}
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
			c.doneCh <- true
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
			c.mux.Lock()
			for _, loco := range c.locomotives {
				loco.sendPackets(c.driver)
			}
			c.mux.Unlock()
			idle.PacketPause()
		}
	}
}
