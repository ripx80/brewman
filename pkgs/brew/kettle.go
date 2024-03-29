package brew

import (
	"errors"
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
	metric   KettleMetric
}

/*KettleMetric for access internal Metrics from outside*/
type KettleMetric struct {
	Temp     float64
	Heater   bool
	Agitator bool
	Fail     int
}

/*CancelErr error if process not finish correctly*/
const CancelErr = "cancel"

/*Metric returns current Metrics from Kettle*/
func (k *Kettle) Metric() KettleMetric {
	// workaround if temp is not init, when no step is currently running
	if k.metric.Temp == 0 {
		k.metric.Temp, _ = k.Temp.Get()
	}
	k.metric.Heater = k.Heater.State()
	k.metric.Agitator = k.Agitator.State()
	return k.metric
}

/*
On turns the Agitator and the Heater on if available
*/
func (k *Kettle) On() {
	k.Agitator.On()
	k.Heater.On()
}

/*
Off turns the Agitator and the Heater of if available
*/
func (k *Kettle) Off() {
	k.Agitator.Off()
	k.Heater.Off()

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
		last     float64
		temp     float64
		failcnt  uint8
		lastfail time.Time
		err      error
	)

	for {
		select {
		case <-stop:
			k.Off()
			return errors.New(CancelErr)
		case <-time.After(1 * time.Second):
			if temp, err = k.TempSet(tempTo); err != nil {
				log.Error("temp sensor failed: %w", err)
			}
			// reset old counter
			if (lastfail != time.Time{}) && time.Now().After(lastfail.Add(time.Second*20)) {
				failcnt = 0
				lastfail = time.Time{}
			}
			if !k.TempCompare(last, temp) {
				failcnt++
				lastfail = time.Now()
			}
			last = temp

			if failcnt >= 10 {
				log.Error("Temperature not increased but the heater is on. Check your hardware setup")
				k.Heater.On() // set it again
				failcnt = 0
				lastfail = time.Time{}
				k.metric.Fail++
			}

			if temp >= tempTo {
				if k.Heater.State() {
					k.Heater.Off()
				}
				return nil
			}

			//use zstate and not state because logrus log in alphabetical order. workaround sry
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
			return errors.New(CancelErr)

		case <-ttl:
			return nil

		case <-time.After(1 * time.Second):
			if temp, err = k.TempSet(tempTo); err != nil {
				log.Error("temp sensor failed: %w", err)
			}

			// if you have 1.5 difference on holding, increase counter
			if k.Heater.State() && temp <= (tempTo-1.5) {
				failcnt++
			}

			// heat protection
			if !k.Heater.State() && temp >= (tempTo+1.0) {
				log.Error("temperature increase but heater is off. check your hardware setup")
				k.Heater.Off() // set it off, high temp is critical, state is incorrect
			}

			if failcnt >= 10 {
				log.Error("temperature not holding. check your hardware setup")
				k.metric.Fail++
				failcnt = 0
			}
			// change this in future
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
	// bufferr, temp sensors sometimes slow
	k.metric.Temp = current
	return
}
