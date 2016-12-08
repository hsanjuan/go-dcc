package dcc

import (
	"bytes"
	"time"
)

// DCC protocol-defined values for reference.
const (
	BitOnePartMinDuration  = 55 * time.Microsecond
	BitOnePartMaxDuration  = 61 * time.Microsecond
	BitZeroPartMinDuration = 95 * time.Microsecond
	BitZeroPartMaxDuration = 9900 * time.Microsecond
	PacketSeparationMin    = 5 * time.Millisecond
	PacketSeparationMax    = 30 * time.Millisecond
	PreambleBitsMin        = 14
)

// Some customizable DCC-related variables.
var (
	BitOnePartDuration  = 55 * time.Microsecond
	BitZeroPartDuration = 100 * time.Microsecond
	PacketSeparation    = 15 * time.Millisecond
	PreambleBits        = 16
)

// HeadlightCompatMode controls if one bit in the speed instruction is
// reserved for headlight. This reduces speed steps from 32 to 16 steps.
var HeadlightCompatMode = false

// DCCPacket represents the unit of information that can be sent to the DCC
// devices in the system. DCCPacket implements the DCC protocol for converting
// the information into DCC-encoded 1 and 0s.
type DCCPacket struct {
	driver  DCCDriver
	address byte
	data    []byte
	ecc     byte

	// encoded holds 1 byte per packet bit because
	// go does not handle single bits well
	encoded []byte
}

// NewPacket returns a new generic DCC Packet.
func NewPacket(d DCCDriver, addr byte, data []byte) *DCCPacket {
	ecc := addr
	for _, i := range data {
		ecc = ecc ^ i
	}

	return &DCCPacket{
		driver:  d,
		address: addr,
		data:    data,
		ecc:     ecc,
	}
}

func NewBaselinePacket(d DCCDriver, addr byte, data []byte) *DCCPacket {
	addr = addr & 0x7F // 0b01111111 last 7 bits
	return NewPacket(d, addr, data)
}

// NewSpeedAndDirectionPacket returns a new baseline DCC packet with speed and
// direction information.
func NewSpeedAndDirectionPacket(d DCCDriver, addr byte, speed byte, dir Direction) *DCCPacket {
	addr = addr & 0x7F // 0b 0111 1111
	if HeadlightCompatMode {
		speed = speed & 0x0F // 4 lower bytes
	} else {
		speed = speed & 0x1F // 5 lower bytes
	}

	dirB := byte(0x1&dir) << 5
	data := (1 << 6) | dirB | speed // 0b01DCSSSS

	return &DCCPacket{
		driver:  d,
		address: addr,
		data:    []byte{data},
		ecc:     addr ^ data,
	}
}

// NewFunctionGroupOnePacket returns an advanced DCC packet which allows to
// control FL,F1-F4 functions. FL is usually associated to the headlights.
func NewFunctionGroupOnePacket(d DCCDriver, addr byte, fl, fl1, fl2, fl3, fl4 bool) *DCCPacket {
	var data, fln, fl1n, fl2n, fl3n, fl4n byte = 0, 0, 0, 0, 0, 0
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

	data = (1 << 7) | fln | fl1n | fl2n | fl3n | fl4n

	return &DCCPacket{
		driver:  d,
		address: addr,
		data:    []byte{data},
		ecc:     addr ^ data,
	}
}

// NewBroadcastResetPacket returns a new broadcast baseline DCC packet which
// makes the decoders erase their volatile memory and return to power up
// state. This stops all locomotives at non-zero speed.
func NewBroadcastResetPacket(d DCCDriver) *DCCPacket {
	return &DCCPacket{
		driver:  d,
		address: 0,
		data:    []byte{0},
		ecc:     0 ^ 0,
	}
}

// NewBroadcastIdlePacket returns a new broadcast baseline DCC packet
// on which decoders perform no action.
func NewBroadcastIdlePacket(d DCCDriver) *DCCPacket {
	return &DCCPacket{
		driver:  d,
		address: 0xFF,
		data:    []byte{0},
		ecc:     0xFF ^ 0,
	}
}

// NewBroadcastStopPacket returns a new broadcast baseline DCC packet which
// tells the decoders to stop all locomotives. If softStop is false, an
// emergency stop will happen by cutting power off the engine.
func NewBroadcastStopPacket(d DCCDriver, dir Direction, softStop bool, ignoreDir bool) *DCCPacket {
	var speed byte = 0
	if !softStop {
		speed = 1
	}

	if ignoreDir {
		speed = speed | (1 << 4)
	}

	dirB := 0x1 & byte(dir)

	data := (1 << 6) | (dirB << 5) | speed

	return &DCCPacket{
		driver:  d,
		address: 0x0,
		data:    []byte{data},
		ecc:     0x0 ^ data,
	}
}

// delayPoll causes a active delay for the specified time
// by actively polling the clock. Unfortunately, for latencies
// under 100us, it is not possible to sleep reliably with
// syscall.Nanosleep().
func delayPoll(d time.Duration) {
	start := time.Now()
	for {
		if time.Since(start) > d {
			return
		}
	}
}

// zero sends a 0 using the DCCDriver
func (p *DCCPacket) zero() {
	p.driver.Low()
	//syscall.Nanosleep(&zeroTs, nil)
	//time.Sleep(BitZeroPartDuration)
	delayPoll(BitZeroPartDuration)
	p.driver.High()
	//syscall.Nanosleep(&zeroTs, nil)
	//time.Sleep(BitZeroPartDuration)
	delayPoll(BitZeroPartDuration)
}

// one sends a 1 using the DCCDriver
func (p *DCCPacket) one() {
	p.driver.Low()
	//syscall.Nanosleep(&oneTs, nil)
	//time.Sleep(BitOnePartDuration)
	delayPoll(BitOnePartDuration)
	p.driver.High()
	//syscall.Nanosleep(&oneTs, nil)
	//time.Sleep(BitOnePartDuration)
	delayPoll(BitOnePartDuration)
}

// PacketPause performs a pause by sleeping
// during the PacketSeparation time.
func (p *DCCPacket) PacketPause() {
	// Not really needed
	p.driver.Low()
	time.Sleep(PacketSeparation)
	p.driver.High()
}

// Send encodes and sends a packet using the DCCDriver associated to it.
func (p *DCCPacket) Send() {
	if p.driver == nil {
		panic("No driver set")
	}
	if p.encoded == nil {
		p.build()
	}

	for _, b := range p.encoded {
		if b == 0 {
			p.zero()
		} else {
			p.one()
		}
	}
}

// By prebuilding packages we ensure more consistent Send() times.
func (p *DCCPacket) build() {
	unpackByte := func(b byte) []byte {
		bs := make([]byte, 8, 8)
		for i := uint8(0); i < 8; i++ {
			bit := (b >> (7 - i)) & 0x1
			bs[i] = bit
		}
		return bs
	}
	var buf bytes.Buffer

	// Preamble
	for i := 0; i < PreambleBits; i++ {
		buf.WriteByte(1)
	}

	// Packet start bit
	buf.WriteByte(0)

	// Address
	buf.Write(unpackByte(p.address))

	// Data
	buf.WriteByte(0) // Data start bit
	for _, d := range p.data {
		buf.Write(unpackByte(d))
		buf.WriteByte(0) // Data start bit
	}
	buf.Write(unpackByte(p.ecc))
	// Packet end
	buf.WriteByte(1)
	p.encoded = buf.Bytes()
}

func (p *DCCPacket) String() string {
	if p.encoded == nil {
		p.build()
	}
	var str string
	for _, b := range p.encoded {
		if b == 0 {
			str += "0"
		} else if b == 1 {
			str += "1"
		} else {
			panic("bad encoding")
		}
	}
	return str
}
