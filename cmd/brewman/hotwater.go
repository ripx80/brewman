package main

import (
	log "github.com/ripx80/brave/log/logger"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/recipe"
)

// remove config file from here!
func Hotwater(configFile *config.Config, stop chan struct{}) error {
	var err error
	kettle := &brew.Kettle{}
	if err = Init(kettle, configFile.Hotwater); err != nil {
		return err
	}
	recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
	if err != nil {
		return err
	}

	log.Infof("using recipe: ", recipe.Global.Name)
	log.Infof("main water: %f -->  grouting: %f", recipe.Water.MainCast, recipe.Water.Grouting)

	if kettle.Agitator != nil && !kettle.Agitator.State() {
		kettle.Agitator.On()
	}

	if err := kettle.TempUp(stop, configFile.Global.HotwaterTemperatur); err != nil {
		return err
	}

	return kettle.TempHold(stop, configFile.Global.HotwaterTemperatur, 0)
}
