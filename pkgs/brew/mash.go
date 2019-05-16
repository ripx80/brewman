package brew

import (
	"fmt"
	"time"
)

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type Kettle struct {
	Temp     TempSensor
	Heater   Control
	Agitator Control
	//Recipe   recipe.RecipeMash
}

//it prints output.. must be implemented in main
func (k *Kettle) ToTemp(temp float64) error {

	var err error
	var current float64

	for {
		if current, err = k.Temp.Get(); err != nil {
			break
		}

		if current >= temp {
			break
		}
		fmt.Printf("%f --> %f On: %t\n", current, temp, k.Heater.State())
		if current < temp {
			if !k.Heater.State() {
				k.Heater.On()
			}
		}
		// sleep interval as config value
		time.Sleep(2 * time.Second)
	}
	k.Heater.Off()
	return err
}

func (k *Kettle) HoldTemp(temp float64, holdTime time.Duration) error {
	var current float64
	var err error
	stop := time.After(holdTime)
	for {
		select {
		case <-stop:
			//finish success
			return nil
		case <-time.After(1 * time.Second):
			if current, err = k.Temp.Get(); err != nil {
				return err
			}
			if current < temp && !k.Heater.State() {
				// log.Debug("Heater Off")
				k.Heater.On()
			}
			if current > temp && k.Heater.State() {
				// log.Debug("Heater Off")
				k.Heater.Off()
			}
			fmt.Printf("Hold: %f --> %f State: %t\n", current, temp, k.Heater.State())
		}
	}
}
