package brew

import (
	"fmt"
	"time"

	log "github.com/ripx80/brave/log/logger"
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
TempCompare checks the temp and sate of the kettle
*/
func (k *Kettle) TempCompare(last float64, temp float64) bool {
	return (k.Heater.State() && temp > (last)) || (!k.Heater.State() && temp < (last))
}

/*
mabye do it not in the lib or in a other place to keep the lib simple
*/

/*
TempUp control the Heater to increase to a given temperature
Its a blocking function which you can stop with the stop channel
*/
func (k *Kettle) TempUp(stop chan struct{}, tempTo float64) error {

	var (
		last    float64
		temp    float64
		failcnt uint8
		err     error
	)

	for {
		select {
		case <-stop:
			k.Off()
			return nil
		case <-time.After(1 * time.Second):
			if temp, err = k.TempSet(tempTo); err != nil {
				return err
			}
			if !k.TempCompare(last, temp) {
				failcnt++
			}
			last = temp

			if failcnt >= 6 {
				log.Error("Temperature not increased but the heater is on. Check your hardware setup")
				failcnt = 0
			}

			if temp >= tempTo {
				if k.Heater.State() {
					k.Heater.Off()
				}
				return nil
			}
			// use zstate and not state because logrus log in alphabetical order. workaround sry
			log.WithFields(log.Fields{
				"temperatur":  fmt.Sprintf("%0.2f", temp),
				"destination": tempTo,
				"zstate":      k.Heater.State(),
				"fail":        failcnt,
			}).Info("increase temperatur")
		}
	}
}

/*
TempHold control the Heater to hold a given temperature. You can set a duration or 0 (unlimited)
Its a blocking function which you can stop with the stop channel
*/
func (k *Kettle) TempHold(stop chan struct{}, tempTo float64, timeout time.Duration) error {
	var (
		temp    float64
		failcnt uint8
		err     error
	)
	ttl := make(<-chan time.Time, 1) // placeholder for timer, 0 run forever
	if timeout > 0 {
		ttl = time.After(timeout)
	}
	for {
		select {
		case <-stop:
			k.Off()
			return nil

		case <-ttl:
			return nil

		case <-time.After(1 * time.Second):
			if temp, err = k.TempSet(tempTo); err != nil {
				return err
			}

			// if you have 1.5 difference on holding, increase counter
			if k.Heater.State() && temp <= (tempTo-1.5) {
				failcnt++
			}

			if failcnt >= 3 {
				log.Error("temperature not holding but the heater is on. check your hardware setup")
				failcnt = 0
			}

			log.WithFields(log.Fields{
				"temperatur":  fmt.Sprintf("%0.2f", temp),
				"destination": tempTo,
				"zstate":      k.Heater.State(),
				"fail":        failcnt,
			}).Info("holding temperatur")
		}
	}
}

/*
TempSet check the state of the Heater and turn off/on related to the given temp
*/
func (k *Kettle) TempSet(temp float64) (current float64, err error) {
	if current, err = k.Temp.Get(); err != nil {
		return 0, err
	}
	if current < temp && !k.Heater.State() {
		log.Debug("Heater On")
		k.Heater.On()
	}
	if current > temp && k.Heater.State() {
		log.Debug("Heater Off")
		k.Heater.Off()
	}
	return
}
