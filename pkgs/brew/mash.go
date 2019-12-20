package brew

import (
	"fmt"
	"os"
	"time"

	"github.com/ripx80/brewman/config"
	log "github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

type Kettle struct {
	Temp     TempSensor
	Heater   Control
	Agitator Control
}

/*
supported sensors
	dummy (dev mode)
	GPIO SSR (control over gpio pins)
	GPIO Temp (ds18b20)
*/
func (k *Kettle) Init(kettleConfig config.PodConfig) error {

	if kettleConfig == (config.PodConfig{}) {
		return fmt.Errorf("no podconfig in config file. you must have a podconfig to mash/hotwater/cooking")
	}
	_, err := host.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize periph: %v", err)
	}

	var p gpio.PinIO

	if kettleConfig.Control.Device == "dummy" {
		k.Heater = &SSRDummy{}
	} else {
		p = gpioreg.ByName(kettleConfig.Control.Address)
		if p == nil {
			return fmt.Errorf("failed to find heater pin: %s", kettleConfig.Control.Address)
		}
		k.Heater = &SSR{Pin: p}
	}

	// Agiator
	switch kettleConfig.Agiator.Device {
	case "dummy":
		k.Agitator = &SSRDummy{}
	case "gpio":
		p = gpioreg.ByName(kettleConfig.Agiator.Address)
		if p == nil {
			return fmt.Errorf("failed to find agiator pin: %s", kettleConfig.Agiator.Address)
		}
		k.Agitator = &SSR{Pin: p}
	case "":
		if kettleConfig.Agiator.Address != "" {
			return fmt.Errorf("failed setup agiator, device not set: %s", kettleConfig.Agiator.Address)
		}
	default:
		return fmt.Errorf("unsupported agiator device: %s", kettleConfig.Agiator.Device)
	}

	// Temperatur
	switch kettleConfig.Temperatur.Device {
	case "ds18b20":
		if _, err := os.Stat(kettleConfig.Temperatur.Address); os.IsNotExist(err) {
			return fmt.Errorf("path to temp sensor not exists: %s", kettleConfig.Temperatur.Address)
		}
		k.Temp = DS18B20{Name: kettleConfig.Temperatur.Device, Path: kettleConfig.Temperatur.Address}
	case "dummy":
		k.Temp = &TempDummy{Name: "tempdummy", fn: k.Heater.State, Temp: 20}
	case "default":
		return fmt.Errorf("unsupported temp device: %s", kettleConfig.Temperatur.Device)
	}
	return nil
}

func (k *Kettle) On() {
	if k.Agitator != nil && !k.Agitator.State() {
		k.Agitator.On()
	}

	if k.Heater != nil && !k.Heater.State() {
		k.Heater.On()
	}
}

func (k *Kettle) Cleanup() {
	if k.Agitator != nil && !k.Agitator.State() {
		k.Agitator.Off()
	}

	if k.Heater != nil && !k.Heater.State() {
		k.Heater.Off()
	}
}

func (k *Kettle) TempIncreaseTo(stop chan struct{}, tempTo float64) error {
	var (
		temp float64
		err  error
	)
	for {
		select {
		case <-stop:
			if k.Heater.State() {
				k.Heater.Off()
			}
			return nil
		case <-time.After(1 * time.Second):
			if temp, err = k.TempWatch(tempTo); err != nil {
				return err
			}
			if temp >= tempTo {
				if k.Heater.State() {
					k.Heater.Off()
				}
				return nil
			}
			// use data channel
			log.Infof("Increase: %f --> %f State: %t\n", temp, tempTo, k.Heater.State())
		}
	}
	return fmt.Errorf("terminated")
}

func (k *Kettle) TempHolder(stop chan struct{}, tempTo float64, holdTime time.Duration) error {
	var (
		temp float64
		err  error
	)
	timeout := make(<-chan time.Time, 1) // placeholder for timer, 0 run forever
	if holdTime > 0 {
		timeout = time.After(holdTime)
	}
	for {
		select {
		case <-stop:
			if k.Heater.State() {
				k.Heater.Off()
			}
			return nil

		case <-timeout:
			return nil

		case <-time.After(1 * time.Second):
			if temp, err = k.TempWatch(tempTo); err != nil {
				return err
			}
			log.Infof("Hold: %f --> %f State: %t\n", temp, tempTo, k.Heater.State())
		}
	}
	return fmt.Errorf("terminated")
}

func (k *Kettle) TempWatch(temp float64) (current float64, err error) {
	if current, err = k.Temp.Get(); err != nil {
		return 0, err
	}
	if current < temp && !k.Heater.State() {
		log.Debug("Heater Off")
		k.Heater.On()
	}
	if current > temp && k.Heater.State() {
		log.Debug("Heater Off")
		k.Heater.Off()
	}
	return
}
