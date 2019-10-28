package dcc

import "sync"

// CommandRepeat specifies how many times a single
// packet is sent.
var CommandRepeat = 30

// CommandMaxQueue specifies how many commands can
// queue before sending a new command blocks
// the sender
var CommandMaxQueue = 3

// Controller represents a DCC Control Station. The
// controller keeps tracks of the DCC Locomotives and
// is in charge of sending DCC packets continuously to
// the tracks.
type Controller struct {
	names       map[string]int
	locoMux     sync.RWMutex
	locomotives []Locomotive

	driver Driver

	started    bool
	doneCh     chan bool
	shutdownCh chan bool
	commandCh  chan Packet

	idle Packet
	stop Packet
}

// NewController builds a Controller.
func NewController(d Driver) *Controller {
	// Prebuild commonly used packets
	idle := NewBroadcastIdlePacket(d)
	stop := NewBroadcastStopPacket(d, Forward, false, true)
	idle.build()
	stop.build()

	return &Controller{
		driver:      d,
		names:       make(map[string]int),
		locomotives: make([]Locomotive, 0, 256), // No more than 256 addresses
		doneCh:      make(chan bool),
		shutdownCh:  make(chan bool),
		commandCh:   make(chan Packet, CommandMaxQueue),
		idle:        idle,
		stop:        stop,
	}
}

// NewControllerWithConfig builds a new Controller using the
// given configuration.
func NewControllerWithConfig(d Driver, cfg Config) *Controller {
	c := NewController(d)

	for _, loco := range cfg.Locomotives {
		c.AddLoco(Locomotive{
			Name:      loco.Name,
			Address:   loco.Address,
			Speed:     loco.Speed,
			Direction: loco.Direction,
			Fl:        loco.Fl,
		})
	}
	return c
}

// AddLoco adds or updates a DCC device in the controller. The device will
// start receiving packets if the controller is running. Devices are indexed
// by name. So adding a second Locomotive will replace a previous one with the
// same name. If you have changed the properties of the Locomotive, add it
// again for changes to take effect.
func (c *Controller) AddLoco(l Locomotive) {
	l.prebuildPackets(c.driver)
	c.locoMux.Lock()
	{
		if pos, ok := c.names[l.Name]; ok {
			c.locomotives[pos] = l
		} else {
			c.locomotives = append(c.locomotives, l)
			c.names[l.Name] = len(c.locomotives) - 1
		}
	}
	c.locoMux.Unlock()
}

// RmLoco removes a DCC device from the controller. There
// will be no longer packets sent to it.
func (c *Controller) RmLoco(n string) {
	c.locoMux.Lock()
	{
		if pos, ok := c.names[n]; ok {
			c.locomotives = append(c.locomotives[:pos], c.locomotives[pos+1:]...)
			delete(c.names, n)
		}
	}
	c.locoMux.Unlock()
}

// GetLoco retrieves a DCC device by its Name. The boolean is
// true if the Locomotive was found.
func (c *Controller) GetLoco(n string) (Locomotive, bool) {
	var l Locomotive
	var found bool
	c.locoMux.RLock()
	{
		if pos, ok := c.names[n]; ok {
			found = true
			l = c.locomotives[pos]
		}
	}
	c.locoMux.RUnlock()
	return l, found
}

// Locos returns a list of all registered Locomotives.
func (c *Controller) Locos() []Locomotive {
	c.locoMux.RLock()
	locos := make([]Locomotive, len(c.locomotives), len(c.locomotives))
	copy(locos, c.locomotives)
	c.locoMux.RUnlock()
	return locos
}

// Command allows to send a custom Packet to the tracks.
// The packet will be sent CommandRepeat times.
func (c *Controller) Command(p Packet) {
	if c.started {
		c.commandCh <- p
		return
	}
	c.send(p)
}

// Start starts the controller: powers on the tracks
// and starts sending packets on them.
func (c *Controller) Start() {
	c.driver.TracksOn()
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
}

// SendPackets sends control packets to all registered Locomotives. Each
// packet is sent CommandRepeat times. For each locomotive a speend and
// direction and a function packet are sent. SendPackets is performed
// automatically when the controller has been started. When no locomotives are
// registered, an idle packet will be sent.
func (c *Controller) SendPackets() {
	c.locoMux.RLock()
	if len(c.locomotives) == 0 {
		c.send(c.idle)
		c.locoMux.RUnlock()
		return
	}

	for j := range c.locomotives {
		for i := 0; i < CommandRepeat; i++ {
			c.locomotives[j].sendPackets()
		}
	}
	c.locoMux.RUnlock()
}

// send sends packet to the tracks CommandRepeat times.
func (c *Controller) send(p Packet) {
	for i := 0; i < CommandRepeat; i++ {
		p.Send()
	}
}

func (c *Controller) run() {
	for {
		select {
		case <-c.shutdownCh:
			for i := 0; i < CommandRepeat; i++ {
				c.stop.Send()
			}
			c.driver.TracksOff()
			c.doneCh <- true
			return
		case p := <-c.commandCh:
			c.send(p)
		default:
			c.SendPackets()
			// This is essentially sleeping and giving time to
			// other goroutines to jump in.
			c.idle.PacketPause()
		}
	}
}
