package brew

import (
	"fmt"
	"os/exec"

	signal "github.com/ripx80/gpio"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

/*
Control interface to turn off and on and get a state of the control unit
*/
type Control interface {
	On() error
	Off() error
	State() bool
}

type transmitOptions struct {
	PulseLength uint
	GpioPin     uint
	Protocol    int
	BitLength   int
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

/*External implements a external programm to control. no args, fix*/
type External struct {
	Cmd   string
	state bool
}

/*Signal implements control over a gpio pin via 433MHZ signals*/
type Signal struct {
	Pin   uint
	Code  uint64
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
func (d *SSRDummy) State() bool { return d.state }

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

/*
SSRReg register a valid gpio address and return a SSR struct
*/
func SSRReg(address string) (*SSR, error) {
	var p gpio.PinIO
	// if you use more than one periph driver do this only once in a higher level
	_, err := host.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize periph: %v", err)
	}

	p = gpioreg.ByName(address)
	if p == nil {
		return nil, fmt.Errorf("failed to find heater pin: %s", address)
	}
	return &SSR{Pin: p}, nil
}

/*On executes the given external programm with args no output*/
func (e *External) On() error {
	cmd := exec.Command(e.Cmd, "1")
	if err := cmd.Run(); err != nil {
		return err
	}
	e.state = true
	return nil
}

/*Off executes the given external programm with args no output*/
func (e *External) Off() error {
	cmd := exec.Command(e.Cmd, "0")
	if err := cmd.Run(); err != nil {
		return err
	}
	e.state = false
	return nil
}

/*State returns the current state of external cmd*/
func (e *External) State() bool {
	return e.state
}

/*On send on code over the given gpio pin*/
func (s *Signal) On() error {
	t := signal.NewTransmitter(s.Pin)
	err := t.Transmit(s.Code, signal.DefaultProtocol, 330, signal.DefaultBitLength)
	if err != nil {
		return err
	}
	t.Wait()
	s.state = true
	return nil
}

/*Off send off code over the given gpio pin (on code - 1 )*/
func (s *Signal) Off() error {
	t := signal.NewTransmitter(s.Pin)
	err := t.Transmit((s.Code - 1), signal.DefaultProtocol, 330, signal.DefaultBitLength)
	if err != nil {
		return err
	}
	t.Wait()
	s.state = false
	return nil
}

/*State returns the current state of gpio*/
func (s *Signal) State() bool {
	return s.state
}
