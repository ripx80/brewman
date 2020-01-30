package main

import (
	"fmt"
	"time"

	log "github.com/ripx80/brave/log/logger"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/recipe"
)

/*Mash implement the mash program*/
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

	log.WithFields(log.Fields{
		"recipe": recipe.Global.Name,
	}).Info("mash information:")
	fmt.Println(recipe.Mash)

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

		log.WithFields(log.Fields{
			"number":     num,
			"time":       rast.Time,
			"temperatur": rast.Temperatur,
		}).Info("rast")

		if err := kettle.TempUp(stop, rast.Temperatur); err != nil {
			return err
		}

		if err := kettle.TempHold(stop, rast.Temperatur, time.Duration(rast.Time*60)*time.Second); err != nil {
			return err
		}
	}
	return nil
}

/*MashRast can jump to a specific rast*/
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
	rast := recipe.Mash.Rests[num-1]

	log.WithFields(log.Fields{
		"recipe":     recipe.Global.Name,
		"number":     num,
		"time":       rast.Time,
		"temperatur": rast.Temperatur,
	}).Info("jump to rast")

	if err := kettle.TempUp(stop, rast.Temperatur); err != nil {
		return err
	}

	return kettle.TempHold(stop, rast.Temperatur, time.Duration(rast.Time*60)*time.Second)
}
