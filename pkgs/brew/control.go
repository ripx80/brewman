package brew

import (
	"fmt"
	"log"

	"periph.io/x/periph/conn/onewire"
	"periph.io/x/periph/conn/onewire/onewiretest"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/ds18b20"
	"periph.io/x/periph/host"
)

func Mash() error {
	err := Init()
	if err != nil {
		return err
	}
	return nil

}

//rewrite log...

func Init() error {

	state, err := host.Init()
	if err != nil {
		log.Fatalf("failed to initialize periph: %v", err)
	}

	// Prints the loaded driver.
	fmt.Printf("Using drivers:\n")
	for _, driver := range state.Loaded {
		fmt.Printf("- %s\n", driver)
	}
	// Prints the driver that were skipped as irrelevant on the platform.
	fmt.Printf("Drivers skipped:\n")
	for _, failure := range state.Skipped {
		fmt.Printf("- %s: %s\n", failure.D, failure.Err)
	}

	// Having drivers failing to load may not require process termination. It
	// is possible to continue to run in partial failure mode.
	fmt.Printf("Drivers failed to load:\n")
	for _, failure := range state.Failed {
		fmt.Printf("- %s: %v\n", failure.D, failure.Err)
	}

	// Use pins, buses, devices, etc.

	// args: onewire

	//bus := &onewiretest.Playback{}
	ops := []onewiretest.IO{
		// Match ROM + Read Scratchpad (init)
		{
			W: []uint8{0x55, 0x28, 0xac, 0x41, 0xe, 0x7, 0x0, 0x0, 0x74, 0xbe},
			R: []uint8{0xe0, 0x1, 0x0, 0x0, 0x3f, 0xff, 0x10, 0x10, 0x3f},
		},
		// Match ROM + Convert
		{
			W:    []uint8{0x55, 0x28, 0xac, 0x41, 0xe, 0x7, 0x0, 0x0, 0x74, 0x44},
			Pull: true,
		},
		// Match ROM + Read Scratchpad (read temp)
		{
			W: []uint8{0x55, 0x28, 0xac, 0x41, 0xe, 0x7, 0x0, 0x0, 0x74, 0xbe},
			R: []uint8{0xe0, 0x1, 0x0, 0x0, 0x3f, 0xff, 0x10, 0x10, 0x3f},
		},
	}
	var addr onewire.Address = 0x740000070e41ac28
	bus := onewiretest.Playback{Ops: ops}
	dev, err := ds18b20.New(&bus, addr, 10)

	if err != nil {
		fmt.Println("invalid resolution")
	}
	fmt.Println(dev.String())

	// Read the temperature.

	e := physic.Env{}
	if err := dev.Sense(&e); err != nil {
		return (err)
	}

	if err := dev.Halt(); err != nil {
		return (err)
	}
	if err := bus.Close(); err != nil {
		return (err)
	}

	fmt.Println(e.Temperature)

	//gpio low example
	// l := gpio.Low;
	// rpi.P1_33.Out(l)

	return nil
}
