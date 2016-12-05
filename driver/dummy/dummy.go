package dummy

import (
	"fmt"
	"time"
)

type DCCDummy struct {
	lasttick time.Time
}

func (d *DCCDummy) Low() {
	d.lasttick = time.Now()
}

func (d *DCCDummy) High() {
	//	dur := time.Since(d.lasttick)
	//	fmt.Println(" Tick Duration:", dur)
	// if dur < 150*time.Microsecond {
	// 	fmt.Print("1")
	// } else if dur < 300*time.Microsecond {
	// 	fmt.Print("0")
	// } else {
	// 	fmt.Println()
	// }
}

func (d *DCCDummy) TracksOff() {
	fmt.Println("->Tracks off")
}

func (d *DCCDummy) TracksOn() {
	fmt.Println("->Tracks on")
	d.lasttick = time.Now()
}
