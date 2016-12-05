package dcc

// Encoder allows to convert a DCC packet into an electrical signal
// following the DCC specification
type DCCDriver interface {
	// Low sets the output to low state
	Low()
	// High sets the output to high
	High()
	// TracksOn turns the tracks on. The exact procedure is left to the
	// implementation, but tracks should be ready to receive packets from
	// this point.
	TracksOn()
	// TracksOff disables the tracks. The exact procedure is left to the
	// implementation, but tracks should not carry any power and all
	// trains should stop after calling it.
	TracksOff()
}
