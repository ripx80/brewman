package brew

import (
	"fmt"
	"os"
	"time"

	// remove config dep
	"github.com/ripx80/brewman/config"
	log "github.com/sirupsen/logrus"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

/*
Kettle is the pod unit with Temp, Heater and Agitator
*/
type Kettle struct {
	Temp     TempSensor
	Heater   Control
	Agitator Control
}

/*
Init setup a kettle with information from a config file
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

/*
On turns the Agitator and the Heater on if available
*/
func (k *Kettle) On() {
	if k.Agitator != nil && !k.Agitator.State() {
		k.Agitator.On()
	}

	if k.Heater != nil && !k.Heater.State() {
		k.Heater.On()
	}
}

/*
Off turns the Agitator and the Heater of if available
*/
func (k *Kettle) Off() {
	if k.Agitator != nil && !k.Agitator.State() {
		k.Agitator.Off()
	}

	if k.Heater != nil && !k.Heater.State() {
		k.Heater.Off()
	}
}

/*
TempIncreaseTo control the Heater to increase to a given temperature
*/
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
			if temp, err = k.tempWatch(tempTo); err != nil {
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
}

/*
TempHolder control the Heater to hold a given temperature. You can set a duration or 0 (unlimited)
*/
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
			if temp, err = k.tempWatch(tempTo); err != nil {
				return err
			}
			log.Infof("Hold: %f --> %f State: %t\n", temp, tempTo, k.Heater.State())
		}
	}
}

func (k *Kettle) tempWatch(temp float64) (current float64, err error) {
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
