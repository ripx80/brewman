package brew

import (
	"errors"

	"periph.io/x/periph/conn/onewire"
	"periph.io/x/periph/conn/onewire/onewiretest"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/ds18b20"
)

/*
https://periph.io/device/
"periph.io/x/periph/host/rpi" for raspi based on http://pinout.xyz/
func Present() bool if on raspi board

*/
type TempSensor interface {
	Get() (*physic.Temperature, error)
}

type DS18B20 struct {
	Name   string
	Device *ds18b20.Dev
}

func (ds *DS18B20) Init(addr onewire.Address) error {
	//var addr onewire.Address = 0x740000070e41ac28
	var err error
	bus := onewiretest.Playback{}
	ds.Device, err = ds18b20.New(&bus, addr, 10)

	if err != nil {
		return errors.New("invalid resolution")
	}
	return nil
}

func (ds *DS18B20) Get() (*physic.Temperature, error) {
	e := physic.Env{}
	if err := ds.Device.Sense(&e); err != nil {
		return nil, err
	}

	if err := ds.Device.Halt(); err != nil {
		return nil, err
	}

	return &e.Temperature, nil
}

func (ds *DS18B20) InitDummy() error {
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
		return errors.New("invalid resolution")
	}

	ds.Device = dev
	return nil
}

// func (registerSensor(Sensor, func))
// func (getHardwareSensors)

// struct OneWire interface

// struct WaterFlow type
// struct Thermometer type
// struct WaterLevel type
