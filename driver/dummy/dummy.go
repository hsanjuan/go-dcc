package dummy

import (
	"bytes"
	"fmt"
	"time"
)

// GuessBuffer will be used by the dummy driver to print the value of packets
// sent.
var GuessBuffer bytes.Buffer

// ByteOneMax configures how long a DCC encoded 1 lasts. A tick lasting
// under this value will be guessed as 1.
var ByteOneMax = 61 * time.Microsecond

// ByteZeroMax configures how long a DCC encoded 0 lasts. A tick lasting
// under this value but more than ByteOneTickMax will be guessed as 0.
var ByteZeroMax = 9900 * time.Microsecond

// Driver implements a mock DCC driver that records how long an output lasted
// and prints 0 or 1 to GuessBuffer.
type Driver struct {
	lasttick time.Time
}

// Low sets output to Low.
func (d *Driver) Low() {
	d.lasttick = time.Now()
}

// High sets output to High. High completes a cycle and adds "0" or "1" to the
// GuessBuffer depending on how long it lasted.
func (d *Driver) High() {
	dur := time.Since(d.lasttick)
	if dur < ByteOneMax {
		GuessBuffer.WriteString("1")
	} else if dur < ByteZeroMax {
		GuessBuffer.WriteString("0")
	} else {
		GuessBuffer.WriteString("\n")
	}
}

// TracksOff simulates tracks being turned off.
func (d *Driver) TracksOff() {
	fmt.Println("-> Dummy driver: Tracks off")
}

// TracksOn simulates Tracks being turned on.
func (d *Driver) TracksOn() {
	fmt.Println("-> Dummy driver: Tracks on")
	GuessBuffer.Reset()
	d.lasttick = time.Now()
}
