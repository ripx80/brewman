package brew

/*
https://periph.io/device/
"periph.io/x/periph/host/rpi" for raspi based on http://pinout.xyz/
func Present() bool if on raspi board

*/
type Sensor interface {
	Get() (float64, error)
}

type OneWireSensor struct {
	Name   string
	Device string
}

// if a device is not handled driectly by periph we can implement conn interface

func (s OneWireSensor) Get() (float64, error) {
	return 4.0, nil
}

// func (registerSensor(Sensor, func))
// func (getHardwareSensors)

// struct OneWire interface

// struct WaterFlow type
// struct Thermometer type
// struct WaterLevel type
