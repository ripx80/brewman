package brew

import (
	"periph.io/x/periph/conn/gpio"
)

type Control interface {
	On() error
	Off() error
	State() bool
}

type SSR struct {
	Pin gpio.PinIO
}

type SSRDummy struct {
	state bool
}

func (d *SSRDummy) On() error {
	d.state = true
	return nil
}
func (d *SSRDummy) Off() error {
	d.state = false
	return nil
}
func (d SSRDummy) State() bool { return d.state }

func (ssr *SSR) On() error {
	l := gpio.High
	ssr.Pin.Out(l)
	return nil
}

func (ssr *SSR) Off() error {
	l := gpio.Low
	ssr.Pin.Out(l)
	return nil
}

func (ssr *SSR) State() bool {
	return bool(ssr.Pin.Read())
}

// type Periph struct {
// 	State       *periph.State
// 	TempSensors map[string]TempSensor
// 	Controls    map[string]Control
// }
