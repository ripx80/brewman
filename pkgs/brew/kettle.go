package brew

import (
	"time"

	log "github.com/sirupsen/logrus"
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
Watch is a blocking func which checks the temp and sate of the kettle
*/
func (k *Kettle) Watch(stop chan struct{}, tolerance int) error {
	var (
		temp float64
		last float64
		err  error
	)
	failcnt := 0

	for {
		select {
		case <-stop:
			log.Debug("Watcher go exit")
			return nil
		case <-time.After(1 * time.Second):
			if temp, err = k.Temp.Get(); err != nil {
				return err
			}
			if (k.Heater.State() && temp < (last)) || (!k.Heater.State() && temp > (last)) {
				failcnt++
			}

			if failcnt >= tolerance {
				log.Warnf("Heater is on/off but temp increase/decrease, temp: %.2f last: %.2f", temp, last)
				failcnt = 0
			}
			last = temp
		}
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
func (k *Kettle) TempHolder(stop chan struct{}, tempTo float64, timeout time.Duration) error {
	var (
		temp float64
		err  error
	)
	ttl := make(<-chan time.Time, 1) // placeholder for timer, 0 run forever
	if timeout > 0 {
		ttl = time.After(timeout)
	}
	for {
		select {
		case <-stop:
			if k.Heater.State() {
				k.Heater.Off()
			}
			return nil

		case <-ttl:
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
