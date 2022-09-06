package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"

	dcc "github.com/hsanjuan/go-dcc"
	"github.com/hsanjuan/go-dcc/driver/dccpi"
	"github.com/hsanjuan/go-dcc/driver/dummy"
	rpio "github.com/stianeikeland/go-rpio/v4"
)

const description = `
dccpi allows to use Raspberry Pi as a DCC station to control DCC-enabled
model trains and accessories.

Running dccpi starts a dccpi prompt which can be used to
run different commands, like starting and stopping the controller.

A dummy driver will be used if the Raspberry Pi GPIO pins are not accessible,
either because the application is executed on a different platform or because
the user running it does not have the necessary rights. For the last case,
try running it as root.

`

var cmds = map[string]cmd{
	"help": cmd{
		Name:      "help",
		ShortDesc: "Show this help",
		LongDesc:  "Shows help!",
	},
	"power": {
		Name:      "power",
		ShortDesc: "Control track power",
		LongDesc: `
Usage: power <on|off>

"on" will start delivering power to the tracks and sending DCC
packets on them. "off" will remove power from the tracks and
stop sending packets.
`},
	"register": {
		Name:      "register",
		ShortDesc: "Add DCC device",
		LongDesc: `
Usage: register <device_name> <address>

This command allows to add a device so it can be controlled. The
device will start receiving DCC control packets addressed to it.
Note that unregistered devices may still act upon broadcast packets.
`},
	"unregister": {
		Name:      "unregister",
		ShortDesc: "Remove DCC device",
		LongDesc: `
Usage: unregister <device_name>

This command removes a device. The device will no longer receive any
packets addressed to it.
`},
	"save": {
		Name:      "save",
		ShortDesc: "Save current devices in configuration file",
		LongDesc: `
Usage: save

This command stores the current list of registered devices in the
dccpi configuration file. Note that any other contents will be replaced.
`},
	"status": {
		Name:      "status",
		ShortDesc: "Show information about devices",
		LongDesc: `
Usage: status [device_name]

This command prints information on registered DCC devices. When called
without arguments, it will print information on all devices, otherwise
it will print information on the named device.
`},
	"speed": {
		Name:      "speed",
		ShortDesc: "Control locomotive speed",
		LongDesc: `
Usage: speed <device_name> <speed>

This command sets the speed of the given device. The device will
receive speed-and-direction packets with the given value.
`},
	"direction": {
		Name:      "direction",
		ShortDesc: "Control locomotive direction",
		LongDesc: `
Usage: direction <device_name> <backward|forward|reverse>

This command sets the direction of a given device.
`},
	// 	"estop": {
	// 		Name:      "estop",
	// 		ShortDesc: "Emergency-stop all locomotives",
	// 		LongDesc: `
	// Usage: estop

	// This commands sends a broadcast command which asks DCC devices to
	// cut the power from all locomotives, causing their immediate stop.
	// `},
	"fl": {
		Name:      "fl",
		ShortDesc: "Control the headlight of a locomotive",
		LongDesc: `
Usage: fl <device_name> <on|off>

This command allows to control the FL function of a locomotive, usually
associated with the headlight.
`},
	"exit": {
		Name:      "exit",
		ShortDesc: "Exit from dccpi",
		LongDesc: `
Usage: exit

This command quits dccpi. Tracks are powered off before exiting.
`},
}

// DefaultConfigPath specifies where to read the configuration from
// if no alternative is provided. init() sets it it to ~/.dccpi
var DefaultConfigPath = ""

// The dccpi line prompt
const Prompt = "dccpi> "

// Command line flags
var (
	configFlag    string
	signalPinFlag uint
	brakePinFlag  uint
)

type cmd struct {
	Name      string
	ShortDesc string
	LongDesc  string
}

type repl struct {
	signalCh chan os.Signal
	doneCh   chan struct{}
	ctrl     *dcc.Controller
	driver   dcc.Driver
}

func perr(f string) {
	fmt.Fprintf(os.Stderr, f+"\n")
}

func check(e error) {
	if e != nil {
		perr(e.Error())
		os.Exit(1)
	}
}

func init() {
	usr, _ := user.Current()
	DefaultConfigPath = filepath.Join(usr.HomeDir, ".dccpi")
	flag.Usage = func() {
		perr("Usage: dccpi [options]")
		perr(description)
		perr("Options:")
		flag.PrintDefaults()

	}

	flag.StringVar(&configFlag, "config", DefaultConfigPath,
		"location of a dccpi configuration file")
	flag.UintVar(&signalPinFlag, "signalPin", uint(dccpi.SignalGPIO),
		"GPIO Pin to use for the DCC signal")
	flag.UintVar(&brakePinFlag, "brakePin", uint(dccpi.BrakeGPIO),
		"GPIO Pin to use for the Brake signal (cuts power from tracks")
	flag.Parse()
}

func main() {
	cfg, err := dcc.LoadConfig(configFlag)
	if err != nil {
		perr("Error: cannot load configuration. Using empty one.")
		cfg = &dcc.Config{}
	}

	dccpi.BrakeGPIO = rpio.Pin(brakePinFlag)
	dccpi.SignalGPIO = rpio.Pin(signalPinFlag)

	var dpi dcc.Driver
	dpi, err = dccpi.NewDCCPi()
	if err != nil {
		perr("Error: DCCPi no available. Using dummy driver.")
		dpi = &dummy.DCCDummy{}
	}

	ctrl := dcc.NewControllerWithConfig(dpi, cfg)

	r := &repl{
		signalCh: make(chan os.Signal, 1),
		doneCh:   make(chan struct{}),
		ctrl:     ctrl,
		driver:   dpi,
	}

	signal.Notify(r.signalCh, os.Interrupt)

	go func() {
		<-r.signalCh
		r.shutdown()
	}()

	go r.run()

	<-r.doneCh
	os.Exit(0)
}

func printPrompt() {
	fmt.Printf("%s", Prompt)
}

func (r *repl) shutdown() {
	fmt.Println()
	r.ctrl.Stop()
	fmt.Println("Tracks powered off")
	close(r.doneCh)
}

// run the read-eval-print-loop for dccpi
func (r *repl) run() {
	wrongArgs := func(c string) {
		perr("Error: Wrong command syntax:")
		fmt.Println(cmds[c].LongDesc)
	}
	notReg := func() {
		perr("Error: device not registered")
	}

	for {
		var cmd, arg1, arg2 string
		printPrompt()
		i, _ := fmt.Scanln(&cmd, &arg1, &arg2)
		if i == 0 {
			continue
		}

		switch cmd {
		case "help":
			if arg1 == "" {
				fmt.Println()
				fmt.Println("Available commands (use \"help <command>\" for information):")
				fmt.Println()
				for k, v := range cmds {
					fmt.Printf("%s - %s\n", k, v.ShortDesc)
				}
				fmt.Println()
			} else {
				h, ok := cmds[arg1]
				if !ok {
					perr("Command does not exist.")
				} else {
					fmt.Println(h.LongDesc)
				}
			}
		case "exit":
			r.shutdown()
			return
		case "power":
			if i != 2 {
				wrongArgs(cmd)
				break
			}
			switch arg1 {
			case "on":
				r.ctrl.Start()
			case "off":
				r.ctrl.Stop()
			default:
				wrongArgs(cmd)
			}
		case "register":
			if i != 3 {
				wrongArgs(cmd)
				break
			}
			n, err := strconv.ParseUint(arg2, 10, 8)
			if err != nil {
				perr("Error: wrong DCC address: " + err.Error())
			}
			l := &dcc.Locomotive{
				Name:    arg1,
				Address: uint8(n),
			}
			r.ctrl.AddLoco(l)
		case "unregister":
			if i != 2 {
				wrongArgs(cmd)
				break
			}
			l, ok := r.ctrl.GetLoco(arg1)
			if !ok {
				notReg()
				break
			}
			r.ctrl.RmLoco(l)
		case "status":
			if i > 2 {
				wrongArgs(cmd)
				break
			}
			if i == 2 {
				l, ok := r.ctrl.GetLoco(arg1)
				if !ok {
					notReg()
					break
				}
				fmt.Println(l.String())
			} else {
				locos := r.ctrl.Locos()
				for _, l := range locos {
					fmt.Println(l.String())
				}
			}
		case "speed":
			if i != 3 {
				wrongArgs(cmd)
				break
			}
			l, ok := r.ctrl.GetLoco(arg1)
			if !ok {
				notReg()
				break
			}
			n, err := strconv.ParseUint(arg2, 10, 8)
			if err != nil {
				perr("Wrong speed value: " + err.Error())
			}
			l.Speed = uint8(n)
			l.Apply()
		case "direction":
			if i != 3 {
				wrongArgs(cmd)
				break
			}
			l, ok := r.ctrl.GetLoco(arg1)
			if !ok {
				notReg()
				break
			}
			switch arg2 {
			case "forward":
				l.Direction = dcc.Forward
				l.Apply()
			case "backward":
				l.Direction = dcc.Backward
				l.Apply()
			case "reverse":
				l.Direction = (l.Direction + 1%2)
				l.Apply()
			default:
				wrongArgs(cmd)
			}
		// case "estop":
		// 	if i != 1 {
		// 		wrongArgs(cmd)
		// 		break
		// 	}
		// 	estop := dcc.NewBroadcastStopPacket(r.driver, dcc.Forward, false, true)
		// 	r.ctrl.Command(estop)
		case "fl":
			if i != 3 {
				wrongArgs(cmd)
				break
			}
			l, ok := r.ctrl.GetLoco(arg1)
			if !ok {
				notReg()
				break
			}
			switch arg2 {
			case "on":
				l.Fl = true
				l.Apply()
			case "off":
				l.Fl = false
				l.Apply()
			default:
				wrongArgs(cmd)
			}
		case "save":
			locos := r.ctrl.Locos()
			cfg := &dcc.Config{
				Locomotives: locos,
			}
			err := cfg.Save(configFlag)
			if err != nil {
				perr("Error: saving configuration: " + err.Error())
				break
			}
			fmt.Println("Configuration saved to", configFlag)
		default:
			l, ok := r.ctrl.GetLoco(cmd)
			if !ok {
				perr("Command not available")
				break
			}
			fmt.Println(l.String())
		}
	}
}
