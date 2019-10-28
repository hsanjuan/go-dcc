package dcc

import "time"

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

// Packet represents the unit of information that can be sent to the DCC
// devices in the system. Packet implements the DCC protocol for converting
// the information into DCC-encoded 1 and 0s.
type Packet struct {
	driver  Driver
	address byte
	data    []byte
	ecc     byte

	// encoded holds an int64 (time.Duration) for each
	// bit in a packet. It is an efficient representation
	// to save extra function calls and IFs when sending
	encoded []time.Duration
}

// NewPacket returns a new generic DCC Packet.
func NewPacket(d Driver, addr byte, data []byte) Packet {
	ecc := addr
	for _, i := range data {
		ecc = ecc ^ i
	}

	return Packet{
		driver:  d,
		address: addr,
		data:    data,
		ecc:     ecc,
	}
}

// NewBaselinePacket returns a new generic baseline packet.
// Baseline packets are different because they use a 128 address
// space. Therefore the address is forced to start with bit 0.
func NewBaselinePacket(d Driver, addr byte, data []byte) Packet {
	addr = addr & 0x7F // 0b01111111 last 7 bits
	return NewPacket(d, addr, data)
}

// NewSpeedAndDirectionPacket returns a new baseline DCC packet with speed and
// direction information.
func NewSpeedAndDirectionPacket(d Driver, addr byte, speed byte, dir Direction) Packet {
	addr = addr & 0x7F // 0b 0111 1111
	if HeadlightCompatMode {
		speed = speed & 0x0F // 4 lower bytes
	} else {
		speed = speed & 0x1F // 5 lower bytes
	}

	dirB := byte(0x1&dir) << 5
	data := (1 << 6) | dirB | speed // 0b01DCSSSS

	return Packet{
		driver:  d,
		address: addr,
		data:    []byte{data},
		ecc:     addr ^ data,
	}
}

// NewFunctionGroupOnePacket returns an advanced DCC packet which allows to
// control FL,F1-F4 functions. FL is usually associated to the headlights.
func NewFunctionGroupOnePacket(d Driver, addr byte, fl, fl1, fl2, fl3, fl4 bool) Packet {
	var data, fln, fl1n, fl2n, fl3n, fl4n byte
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

	return Packet{
		driver:  d,
		address: addr,
		data:    []byte{data},
		ecc:     addr ^ data,
	}
}

// NewBroadcastResetPacket returns a new broadcast baseline DCC packet which
// makes the decoders erase their volatile memory and return to power up
// state. This stops all locomotives at non-zero speed.
func NewBroadcastResetPacket(d Driver) Packet {
	return Packet{
		driver:  d,
		address: 0,
		data:    []byte{0},
		ecc:     0 ^ 0,
	}
}

// NewBroadcastIdlePacket returns a new broadcast baseline DCC packet
// on which decoders perform no action.
func NewBroadcastIdlePacket(d Driver) Packet {
	return Packet{
		driver:  d,
		address: 0xFF,
		data:    []byte{0},
		ecc:     0xFF ^ 0,
	}
}

// NewBroadcastStopPacket returns a new broadcast baseline DCC packet which
// tells the decoders to stop all locomotives. If softStop is false, an
// emergency stop will happen by cutting power off the engine.
func NewBroadcastStopPacket(d Driver, dir Direction, softStop bool, ignoreDir bool) Packet {
	var speed byte
	if !softStop {
		speed = 1
	}

	if ignoreDir {
		speed = speed | (1 << 4)
	}

	dirB := 0x1 & byte(dir)

	data := (1 << 6) | (dirB << 5) | speed

	return Packet{
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
func delayPoll(now time.Time, d time.Duration) {
	for {
		if time.Since(now) > d {
			return
		}
	}
}

// PacketPause performs a pause by sleeping
// during the PacketSeparation time.
func (p Packet) PacketPause() {
	if p.driver == nil {
		return // NOP
	}
	// Not really needed
	p.driver.Low()
	time.Sleep(PacketSeparation)
	p.driver.High()
}

// Send encodes and sends a packet using the Driver associated to it.
func (p Packet) Send() {
	if p.driver == nil {
		return // NOP
	}
	if p.encoded == nil {
		p.build()
	}

	// The way p.encoded is we reduce function calls and ifs
	for _, b := range p.encoded {
		p.driver.Low()
		now := time.Now()
		delayPoll(now, b)
		p.driver.High()
		now = time.Now()
		delayPoll(now, b)
	}
}

// Length returns the length of the DCC-encoded representation
// of a packet.
func (p Packet) Length() int {
	l := 0
	l += PreambleBits // Preamble
	l += 1            // Packet start
	l += 8            // Address byte
	for i := 0; i < len(p.data); i++ {
		l += 1 // Data start
		l += 8 // Data byte
	}
	l += 1 // ECC start
	l += 8 // ECC byte
	l += 1 // Packet end
	return l
}

// given a byte, returns the array of durations for every bit.
func unpackByte(b byte) [8]time.Duration {
	var bs [8]time.Duration
	for i := uint8(0); i < 8; i++ {
		bitone := (b >> (7 - i)) & 0x1
		if bitone == 0 {
			bs[i] = BitZeroPartDuration
		} else {
			bs[i] = BitOnePartDuration
		}
	}
	return bs
}

// By prebuilding packages we ensure more consistent Send() times.
func (p *Packet) build() {
	enc := make([]time.Duration, 0, p.Length())

	// Preamble
	for i := 0; i < PreambleBits; i++ {
		enc = append(enc, BitOnePartDuration)
	}

	// Packet start bit
	enc = append(enc, BitZeroPartDuration)

	// Address
	databits := unpackByte(p.address)
	enc = append(enc, databits[:]...)

	// Data
	for _, d := range p.data {
		databits = unpackByte(d)
		enc = append(enc, BitZeroPartDuration) // Data start
		enc = append(enc, databits[:]...)      // Data
	}

	// ECC
	databits = unpackByte(p.ecc)
	enc = append(enc, BitZeroPartDuration) // ECC start
	enc = append(enc, databits[:]...)

	// Packet end
	enc = append(enc, BitOnePartDuration)

	p.encoded = enc
}

func (p Packet) String() string {
	if p.encoded == nil {
		p.build()
	}
	var str string
	for _, b := range p.encoded {
		if b == BitZeroPartDuration {
			str += "0"
		} else if b == BitOnePartDuration {
			str += "1"
		} else {
			panic("bad encoding")
		}
	}
	return str
}
