package brew

import (
	"github.com/ripx80/brewman/pkgs/recipe"
)

type Masher struct {
	Temp     TempSensor
	Heater   Control
	Agitator Control
	Recipe   recipe.RecipeMash
	Log      chan string
}

func (m *Masher) Mash(log chan string) error {
	log <- "Start Mashing"
	temp, err := m.Temp.Get()
	if err != nil {
		return err
	}
	log <- temp.String()
	return nil
}
