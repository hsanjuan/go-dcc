package dummy

import (
	"bytes"
	"fmt"
	"time"
)

// GuessBuffer will be used by the dummy driver to
// print the value of packets sent.
var GuessBuffer bytes.Buffer

// ByteOneTickMax configures how long a DCC encoded
// 1 lasts. A tick lasting under this value will be guessed as 1.
var ByteOneMax = 61 * time.Microsecond

// ByteZeroTickMax configures how long a DCC encoded
// 0 lasts. A tick lasting under this value  but more than
// ByteOneTickMax will be guessed as 0.
var ByteZeroMax = 9900 * time.Microsecond

type DCCDummy struct {
	lasttick time.Time
}

func (d *DCCDummy) Low() {
	d.lasttick = time.Now()
}

func (d *DCCDummy) High() {
	dur := time.Since(d.lasttick)
	if dur < ByteOneMax {
		GuessBuffer.WriteString("1")
	} else if dur < ByteZeroMax {
		GuessBuffer.WriteString("0")
	} else {
		GuessBuffer.WriteString("\n")
	}
}

func (d *DCCDummy) TracksOff() {
	fmt.Println("-> Dummy driver: Tracks off")
}

func (d *DCCDummy) TracksOn() {
	fmt.Println("-> Dummy driver: Tracks on")
	GuessBuffer.Reset()
	d.lasttick = time.Now()
}
