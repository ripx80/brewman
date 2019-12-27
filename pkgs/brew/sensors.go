package brew

/*
supported sensors
	dummy (dev mode)
	GPIO SSR (control over gpio pins)
	GPIO Temp (ds18b20)
*/

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

/*
TempSensor interface to Get() data from sensor
*/
type TempSensor interface {
	Get() (float64, error)
}

/*
DS18B20 temperatur sensor with 1-wire protocol
*/
type DS18B20 struct {
	Name string
	Path string
}

/*
TempDummy struct is a dummy sensor that gives you values for dev and testing
*/
type TempDummy struct {
	Name string
	Temp float64
	Fn   func() bool
}

func down(x float64) float64 {
	return x - 3
}

func up(x float64) float64 {
	return x + 3
}

func upDown(x float64) float64 {
	if math.Mod(x, 2.0) > 0 {
		return x + 3.0
	}
	return x - 3.0
}

/*
Get data from Dummy
*/
func (td *TempDummy) Get() (float64, error) {

	if td.Fn() {
		td.Temp = up(td.Temp)
	}

	if !td.Fn() {
		td.Temp = down(td.Temp)
	}

	if td.Temp < 0 {
		return 0, fmt.Errorf("negative value detected")
	}

	return td.Temp, nil
}

/*
Get data from DS18B20
*/
func (ds *DS18B20) Get() (float64, error) {

	data, err := ioutil.ReadFile(ds.Path)
	if err != nil {
		return 0, err
	}

	str := string(data[len(data)-6 : len(data)-1])
	temp, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("canot read from %s: %s", ds.Path, err)
	}
	temp = temp / 1000
	if temp < 0 {
		return 0, fmt.Errorf("negative value detected %s: %s", ds.Path, err)
	}
	return temp, nil
}

/*
DS18B20Reg register a valid Path for DS18B20 and return a DS18B20 struct
*/
func DS18B20Reg(address string) (*DS18B20, error) {
	if _, err := os.Stat(address); os.IsNotExist(err) {
		return nil, fmt.Errorf("path to temp sensor not exists: %s", address)
	}
	return &DS18B20{Name: "ds18b20", Path: address}, nil

}
