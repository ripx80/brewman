package brew

import (
	"periph.io/x/periph/conn/gpio"
)

/*
Control interface to turn off and on and get a state of the control unit
*/
type Control interface {
	On() error
	Off() error
	State() bool
}

/*
SSR struct with gpio pin for a SSR relay
*/
type SSR struct {
	Pin gpio.PinIO
}

/*
SSRDummy is a dummy SSR relay
*/
type SSRDummy struct {
	state bool
}

/*
On set state of the dummy to true
*/
func (d *SSRDummy) On() error {
	d.state = true
	return nil
}

/*
Off set state of the dummy to false
*/
func (d *SSRDummy) Off() error {
	d.state = false
	return nil
}

/*
State get the current state of dummy
*/
func (d SSRDummy) State() bool { return d.state }

/*
On turn ssr on
*/
func (ssr *SSR) On() error {
	l := gpio.High
	ssr.Pin.Out(l)
	return nil
}

/*
Off turn ssr off
*/
func (ssr *SSR) Off() error {
	l := gpio.Low
	ssr.Pin.Out(l)
	return nil
}

/*
State returns the current state of ssr
*/
func (ssr *SSR) State() bool {
	return bool(ssr.Pin.Read())
}
