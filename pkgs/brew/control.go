package brew

import (
	"errors"
	"fmt"

	"periph.io/x/periph"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host"
)

type Control interface {
	On() error
	Off() error
}

type SSR struct {
	Pin gpio.PinIO
}

func (ssr *SSR) On() error {
	fmt.Println("Turn On")
	l := gpio.High
	ssr.Pin.Out(l)
	return nil
}

func (ssr *SSR) Off() error {
	l := gpio.Low
	ssr.Pin.Out(l)
	fmt.Println("Turn Off")
	return nil
}

type Periph struct {
	State       *periph.State
	TempSensors map[string]TempSensor
	Controls    map[string]Control
}

func (p *Periph) Init() error {

	state, err := host.Init()

	if err != nil {
		return errors.New(fmt.Sprintf("failed to initialize periph: %v", err))
	}
	p.State = state
	return nil
}
