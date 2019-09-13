package brew

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

/*
https://periph.io/device/
"periph.io/x/periph/host/rpi" for raspi based on http://pinout.xyz/
func Present() bool if on raspi board

*/
type TempSensor interface {
	Get() (float64, error)
}

type DS18B20 struct {
	Name string
	Path string
}

type TempDummy struct {
	Name string
	Temp float64
	fn   func() bool
}

func Down(x float64) float64 {
	return x - 3
}

func Up(x float64) float64 {
	return x + 3
}

func UpDown(x float64) float64 {
	if math.Mod(x, 2.0) > 0 {
		return x + 3.0
	}
	return x - 3.0
}

func (td *TempDummy) Get() (float64, error) {

	if td.fn() {
		td.Temp = Up(td.Temp)
	}

	if !td.fn() {
		td.Temp = Down(td.Temp)
	}

	if td.Temp < 0 {
		return 0, fmt.Errorf("negative value detected")
	}

	return td.Temp, nil
}

func (ds DS18B20) Get() (float64, error) {

	data, err := ioutil.ReadFile(ds.Path)
	if err != nil {
		return 0, err
	}

	str := string(data[len(data)-6 : len(data)-1])
	temp, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Printf("canot read from %s: %s", ds.Path, err)
		return 0, fmt.Errorf("canot read from %s: %s", ds.Path, err)
	}
	temp = temp / 1000
	if temp < 0 {
		return 0, fmt.Errorf("negative value detected %s: %s", ds.Path, err)
	}
	return temp, nil
}

// func (registerSensor(Sensor, func))
// func (getHardwareSensors)

// struct OneWire interface

// struct WaterFlow type
// struct Thermometer type
// struct WaterLevel type
