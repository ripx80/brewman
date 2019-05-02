package brew

import (
	"fmt"

	"periph.io/x/periph/conn/physic"
)

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}

type Kettle struct {
	Temp     TempSensor
	Heater   Control
	Agitator Control
	//Recipe   recipe.RecipeMash
}

func (k *Kettle) GoToTemp(tempTo physic.Temperature) error {

	current, err := k.Temp.Get()
	if err != nil {
		return err
	}

	fmt.Printf("%f < %f\n", current, tempTo)

	// if current < tempTo {
	// 	fmt.Println(k.Heater.State())
	// 	if !k.Heater.State() {
	// 		k.Heater.On()
	// 	}
	// }
	k.Heater.Off()

	return nil

}

func (k *Kettle) None() {}
