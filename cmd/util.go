package cmd

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
)

// func absolutePath(fp *os.File) (string, error) {
// 	return filepath.Abs(fp.Name())
// }

func getControl(device, address string, required bool) (brew.Control, error) {
	var (
		control brew.Control
		err     error
	)
	switch device {
	case "dummy":
		control = &brew.SSRDummy{}
	case "gpio":
		control, err = brew.SSRReg(address)
		if err != nil {
			return nil, err
		}
	case "signal":
		code, err := strconv.Atoi(address)
		if err != nil {
			return nil, err
		}
		control = &brew.Signal{Pin: 17, Code: uint64(code)}

	case "external":
		_, err = os.Stat(address)
		if err != nil {
			return nil, err
		}
		control = &brew.External{Cmd: address}
	case "":
		// can be null no error
		if !required && address != "" {
			return nil, fmt.Errorf("failed setup agiator, device not set: %s", address)
		}
		control = &brew.SSRDummy{}

	default:
		return nil, fmt.Errorf("unsupported control device: %s", device)
	}
	return control, nil
}

func getTempSensor(device, address string, state func() bool) (brew.TempSensor, error) {
	var (
		sensor brew.TempSensor
		err    error
	)
	switch device {
	case "ds18b20":
		sensor, err = brew.DS18B20Reg(address)
		if err != nil {
			return nil, err
		}
	case "dummy":
		sensor = &brew.TempDummy{Name: "tempdummy", Fn: state, Temp: 20}
	case "default":
		return nil, fmt.Errorf("unsupported temp device: %s", device)
	}
	return sensor, nil
}

func getKettle(kconf config.PodConfig) (*brew.Kettle, error) {
	var err error
	k := &brew.Kettle{}
	if kconf == (config.PodConfig{}) {
		return nil, fmt.Errorf("no cod config in config file. you must have a cod config to mash/hotwater/cooking")
	}

	// Control Unit
	k.Heater, err = getControl(kconf.Control.Device, kconf.Control.Address, true)
	if err != nil {
		return nil, err
	}

	// Control Unit
	k.Agitator, err = getControl(kconf.Agiator.Device, kconf.Agiator.Address, false)
	if err != nil {
		return nil, err
	}

	// Temperatur
	k.Temp, err = getTempSensor(kconf.Temperatur.Device, kconf.Temperatur.Address, k.Heater.State)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func goExit(signals chan os.Signal) {
	signals <- syscall.SIGINT // stops all threats and do a cleanup
	select {}
}
