package brew

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heptio/workgroup"
	"github.com/ripx80/brewman/config"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type Kettle struct {
	Temp     TempSensor
	Heater   Control
	Agitator Control
}

func (k *Kettle) Init(kettleConfig config.PodConfig) error {

	if kettleConfig == (config.PodConfig{}) {
		return fmt.Errorf("No Masher in config file. You must have a masher to mash")
	}

	_, err := host.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize periph: %v", err)
	}

	//Heater
	p := gpioreg.ByName(kettleConfig.Control)
	if p == nil {
		return fmt.Errorf("Failed to find Heater Pin: %s", kettleConfig.Agiator)
	}
	k.Heater = &SSR{Pin: p}

	// Agiator
	if kettleConfig.Agiator != "" {
		p = gpioreg.ByName(kettleConfig.Agiator)
		if p == nil {
			return fmt.Errorf("Failed to find Agiator Pin: %s", kettleConfig.Agiator)
		}
		k.Agitator = &SSR{Pin: p}
	}

	// Temperatur
	k.Temp = DS18B20{Name: kettleConfig.Temperatur.Device, Path: kettleConfig.Temperatur.Address}
	if err != nil {
		return fmt.Errorf("Failed to register Temp Sensor: %s", err)
	}
	return nil
}

func (k *Kettle) TempIncreaseTo(tempTo float64) error {
	var g workgroup.Group
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	g.Add(func(stop <-chan struct{}) error {
		select {
		case <-signals:
		case <-stop:
		}
		return fmt.Errorf("stopped")
	})

	// worker thread
	g.Add(func(stop <-chan struct{}) error {
		var (
			temp float64
			err  error
		)
		for {
			select {
			case <-stop:
				fmt.Println("cleaning up")
				if k.Heater.State() {
					k.Heater.Off()
				}
				return fmt.Errorf("terminated")
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
				//implement global log interface
				fmt.Printf("Increase: %f --> %f State: %t\n", temp, tempTo, k.Heater.State())
			}
		}
	})
	return g.Run()
}

func (k *Kettle) TempHolder(tempTo float64, holdTime time.Duration) error {
	var g workgroup.Group
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	timeout := make(<-chan time.Time, 1) // placeholder for timer
	if holdTime > 0 {
		timeout = time.After(holdTime)
	}

	g.Add(func(stop <-chan struct{}) error {
		select {
		case <-timeout:
			return nil // this works for hold time.. not for realy timeouts
		case <-signals:
		case <-stop:
		}
		return fmt.Errorf("stopped")
	})

	// worker thread
	g.Add(func(stop <-chan struct{}) error {
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
				return fmt.Errorf("stopped")
			case <-time.After(1 * time.Second):
				if temp, err = k.TempWatch(tempTo); err != nil {
					return err
				}
				//implement global log interface
				fmt.Printf("Hold: %f --> %f State: %t\n", temp, tempTo, k.Heater.State())
			}
		}
	})

	return g.Run()

}

func (k *Kettle) TempWatch(temp float64) (current float64, err error) {
	if current, err = k.Temp.Get(); err != nil {
		return 0, err
	}
	if current < temp && !k.Heater.State() {
		// log.Debug("Heater Off")
		k.Heater.On()
	}
	if current > temp && k.Heater.State() {
		// log.Debug("Heater Off")
		k.Heater.Off()
	}
	return
}

// func (k *Kettle) HoldTempDuration(stop <-chan time.Time, temp float64) error {
// 	var current float64
// 	var err error
// 	for {
// 		select {
// 		case <-stop:
// 			if k.Heater.State() {
// 				k.Heater.Off()
// 			}

// 			return nil
// 		case <-time.After(1 * time.Second):
// 			if current, err = k.Temp.Get(); err != nil {
// 				return err
// 			}
// 			if current < temp && !k.Heater.State() {
// 				// log.Debug("Heater Off")
// 				k.Heater.On()
// 			}
// 			if current > temp && k.Heater.State() {
// 				// log.Debug("Heater Off")
// 				k.Heater.Off()
// 			}
// 			fmt.Printf("Hold: %f --> %f State: %t\n", current, temp, k.Heater.State())
// 		}
// 	}
// }

// func (k *Kettle) HoldTemp(done chan struct{}, temp float64) error {
// 	var (
// 		err     error
// 		current float64
// 	)
// 	for {
// 		select {
// 		case <-done:
// 			if k.Heater.State() {
// 				k.Heater.Off()
// 			}
// 			return nil
// 		case <-time.After(1 * time.Second):
// 			if current, err = k.Temp.Get(); err != nil {
// 				return err
// 			}
// 			if current < temp && !k.Heater.State() {
// 				// log.Debug("Heater Off")
// 				k.Heater.On()
// 			}
// 			if current > temp && k.Heater.State() {
// 				// log.Debug("Heater Off")
// 				k.Heater.Off()
// 			}
// 		}
// 	}
// }
