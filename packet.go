package dcc

import (
	"log"
	"syscall"
	"time"
)

// DCC protocol-defined values
const (
	BitOnePartMinDuration  = 55 * time.Microsecond
	BitOnePartMaxDuration  = 61 * time.Microsecond
	BitZeroPartMinDuration = 95 * time.Microsecond
	BitZeroPartMaxDuration = 9900 * time.Microsecond
	PacketSeparationMin    = 5 * time.Millisecond
	PacketSeparationMax    = 30 * time.Millisecond
	PreambleBitsMin        = 14
)

// Direction
const (
	Backward Direction = 0
	Forward  Direction = 1
)

// Some customizable DCC-related variables
var (
	BitOnePartDuration  = 55 * time.Microsecond
	BitZeroPartDuration = 200 * time.Microsecond
	zeroTs              syscall.Timespec
	oneTs               syscall.Timespec

	PacketSeparation = 15 * time.Millisecond
	PreambleBits     = 16
)

// HeadlightCompatMode controls if one bit in the speed instruction is
// reserved for headlight. This reduces speed steps from 32 to 16
var HeadlightCompatMode = false

var timeval syscall.Timeval

func init() {
	log.Println("Adjusting timers")
	var delay0 time.Duration = 0
	var delay1 time.Duration = 0

	t := time.Now()
	delayPoll(BitOnePartDuration)
	delay1 = time.Since(t) - BitOnePartDuration

	t = time.Now()
	delayPoll(BitZeroPartDuration)
	delay0 = time.Since(t) - BitZeroPartDuration

	if delay1 > BitOnePartDuration {
		delay1 = BitOnePartDuration - 10
	}

	if delay0 > BitOnePartDuration {
		delay0 = BitOnePartDuration - 10
	}

	//BitOnePartDuration -= delay1
	//BitZeroPartDuration -= delay0

	log.Printf("New times: %d, %d.\n", BitOnePartDuration, BitZeroPartDuration)
}

type Direction byte

type DCCPacket struct {
	driver       DCCDriver
	address      byte
	instructions []byte
	error        byte
}

func NewPacket(d DCCDriver, addr byte, ins []byte) *DCCPacket {
	addr = addr & 0x7F // 0b 0111 1111y
	error := addr
	for _, i := range ins {
		error = error ^ i
	}

	return &DCCPacket{
		driver:       d,
		address:      addr,
		instructions: ins,
		error:        error,
	}
}

func NewSpeedAndDirectionPacket(d DCCDriver, addr byte, speed byte, dir Direction) *DCCPacket {
	addr = addr & 0x7F // 0b 0111 1111
	if HeadlightCompatMode {
		speed = speed & 0x0F // 4 lower bytes
	} else {
		speed = speed & 0x1F // 5 lower bytes
	}

	dirB := byte(0x1&dir) << 5
	ins := (1 << 6) | dirB | speed // 0b01DCSSSS

	return &DCCPacket{
		driver:       d,
		address:      addr,
		instructions: []byte{ins},
		error:        addr ^ ins,
	}
}

func NewFunctionGroupOnePacket(d DCCDriver, addr byte, fl, fl1, fl2, fl3, fl4 bool) *DCCPacket {
	var ins, fln, fl1n, fl2n, fl3n, fl4n byte = 0, 0, 0, 0, 0, 0
	if fl {
		fln = 1 << 4
	}
	if fl1 {
		fl1n = 1
	}
	if fl2 {
		fl2n = 1 << 1
	}
	if fl3 {
		fl3n = 1 << 2
	}
	if fl4 {
		fl4n = 1 << 3
	}

	ins = (1 << 7) | fln | fl1n | fl2n | fl3n | fl4n

	return &DCCPacket{
		driver:       d,
		address:      addr,
		instructions: []byte{ins},
		error:        addr ^ ins,
	}
}

func NewBroadcastResetPacket(d DCCDriver) *DCCPacket {
	return &DCCPacket{
		driver:       d,
		address:      0,
		instructions: []byte{0},
		error:        0 ^ 0,
	}
}

func NewBroadcastIdlePacket(d DCCDriver) *DCCPacket {
	return &DCCPacket{
		driver:       d,
		address:      0xFF,
		instructions: []byte{0},
		error:        0xFF ^ 0,
	}
}

func NewBroadcastStopPacket(d DCCDriver, dir Direction, softStop bool, ignoreDir bool) *DCCPacket {
	var speed byte = 0
	if !softStop {
		speed = 1
	}

	if ignoreDir {
		speed = speed | (1 << 4)
	}

	dirB := 0x1 & byte(dir)

	ins := (1 << 6) | (dirB << 5) | speed

	return &DCCPacket{
		driver:       d,
		address:      0x0,
		instructions: []byte{ins},
		error:        0x0 ^ ins,
	}
}

func delayPoll(d time.Duration) {
	syscall.Gettimeofday(&timeval)
	start := time.Duration(timeval.Usec) * time.Microsecond
	for {
		syscall.Gettimeofday(&timeval)
		gone := time.Duration(timeval.Usec) * time.Microsecond
		if gone-start >= d {
			break
		}
	}
}

func delayPoll2(d time.Duration) {
	start := time.Now()
	for {
		if time.Since(start) > d {
			return
		}
	}
}

func (p *DCCPacket) zero() {
	debug("0")
	p.driver.Low()
	//syscall.Nanosleep(&zeroTs, nil)
	//time.Sleep(BitZeroPartDuration)
	delayPoll2(BitZeroPartDuration)
	p.driver.High()
	//syscall.Nanosleep(&zeroTs, nil)
	//time.Sleep(BitZeroPartDuration)
	delayPoll2(BitZeroPartDuration)
}

func (p *DCCPacket) one() {
	debug("1")
	p.driver.Low()
	//syscall.Nanosleep(&oneTs, nil)
	//time.Sleep(BitOnePartDuration)
	delayPoll2(BitOnePartDuration)
	p.driver.High()
	//syscall.Nanosleep(&oneTs, nil)
	//time.Sleep(BitOnePartDuration)
	delayPoll2(BitOnePartDuration)
}

func (p *DCCPacket) bit(b byte) {
	if b == 1 {
		p.one()
	} else if b == 0 {
		p.zero()
	}
}

func (p *DCCPacket) byte(b byte) {
	for i := uint8(0); i < 8; i++ {
		bit := (b >> (7 - i)) & 0x1
		p.bit(bit)
	}
	debug(" ")
}

func (p *DCCPacket) PacketSpace() {
	// Not really needed
	p.driver.Low()
	time.Sleep(PacketSeparation)
	p.driver.High()
	debug("...\n")
}

func (p *DCCPacket) Send() {
	if p.driver == nil {
		log.Println("No driver set")
		return
	}

	// Preamble
	for i := 0; i < PreambleBits; i++ {
		p.bit(1)
	}
	debug(" ")

	// Packet start bit
	p.bit(0)
	debug(" ")

	// Address
	p.byte(p.address)

	// Data
	p.bit(0) // Data start bit
	debug(" ")
	for _, ins := range p.instructions {
		p.byte(ins)
		p.bit(0) //Data start bit
		debug(" ")

	}
	p.byte(p.error)
	// Packet end
	p.bit(1)
	debug("\n")
}
