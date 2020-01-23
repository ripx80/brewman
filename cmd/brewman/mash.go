package main

import (
	"fmt"
	"time"

	log "github.com/ripx80/brave/log/logger"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/recipe"
)

func Mash(configFile *config.Config, stop chan struct{}) error {
	var err error
	kettle := &brew.Kettle{}
	if err = Init(kettle, configFile.Masher); err != nil {
		return err
	}

	recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
	if err != nil {
		return err
	}

	log.Infof("using recipe: ", recipe.Global.Name)
	log.Infof("mash information: ", recipe.Mash)

	if !confirm("start mashing? <y/n>") {
		return nil
	}

	if kettle.Agitator != nil && !kettle.Agitator.State() {
		kettle.Agitator.On()
	}

	if err := kettle.TempUp(stop, recipe.Mash.InTemperatur); err != nil {
		return err
	}

	if !confirm("malt added? continue? <y/n>") {
		return nil
	}

	for num, rast := range recipe.Mash.Rests {
		log.Infof("Rast %d: Time: %d Temperatur:%f\n", num, rast.Time, rast.Temperatur)

		if err := kettle.TempUp(stop, rast.Temperatur); err != nil {
			return err
		}

		if err := kettle.TempHold(stop, rast.Temperatur, time.Duration(rast.Time*60)*time.Second); err != nil {
			return err
		}
	}
	return nil
}

func MashRast(configFile *config.Config, stop chan struct{}, rastNum int) error {
	num := rastNum // improve
	var err error
	kettle := &brew.Kettle{}
	if err = Init(kettle, configFile.Masher); err != nil {
		return err
	}

	recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
	if err != nil {
		return err
	}

	if num > 8 || num <= 0 {
		return fmt.Errorf("rast number out of range [1-8]")
	}

	if len(recipe.Mash.Rests) < num {
		return fmt.Errorf("rast number not in recipe")
	}

	log.Infof("jump to rast number: %d", num)
	log.Infof("using recipe: ", recipe.Global.Name)

	rast := recipe.Mash.Rests[num-1]
	log.Infof("Rast %d: Time: %d Temperatur: %.2f\n", num, rast.Time, rast.Temperatur)
	if err := kettle.TempUp(stop, rast.Temperatur); err != nil {
		return err
	}

	return kettle.TempHold(stop, rast.Temperatur, time.Duration(rast.Time*60)*time.Second)
}
