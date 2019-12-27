package brew

import (
	log "github.com/sirupsen/logrus"
	"time"
)

/*
Todo:
 - Remove logs in funcs. Do this with data chan or grab data outside
*/

/*
Kettle is the pod unit with Temp, Heater and Agitator
*/
type Kettle struct {
	Temp     TempSensor
	Heater   Control
	Agitator Control
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
