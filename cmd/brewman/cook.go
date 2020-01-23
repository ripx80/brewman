package main

import (
	"time"

	log "github.com/ripx80/brave/log/logger"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/recipe"
)

func Cook(configFile *config.Config, stop chan struct{}) error {
	var err error
	kettle := &brew.Kettle{}
	if err = Init(kettle, configFile.Cooker); err != nil {
		return err
	}

	recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
	if err != nil {
		return err
	}

	log.Infof("using recipe: ", recipe.Global.Name)
	log.Infof("cook information: ", recipe.Cook)

	if !confirm("start cooking? <y/n>") {
		return nil
	}

	if kettle.Agitator != nil && !kettle.Agitator.State() {
		kettle.Agitator.On()
	}

	if err := kettle.TempUp(stop, configFile.Global.CookingTemperatur); err != nil {
		return err
	}

	return kettle.TempHold(stop, configFile.Global.CookingTemperatur, time.Duration(recipe.Cook.Time*60)*time.Second)

}
