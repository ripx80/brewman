package main

import (
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
)

/*ControlOff turn a cattle off*/
func ControlOff(podConfig config.PodConfig) error {
	kettle := &brew.Kettle{}
	if err := Init(kettle, podConfig); err != nil {
		return err
	}
	kettle.Off() // no return value?
	return nil
}

/*ControlOn turn a kettle on*/
func ControlOn(podConfig config.PodConfig) error {
	kettle := &brew.Kettle{}
	if err := Init(kettle, podConfig); err != nil {
		return err
	}
	kettle.On() // no return value?
	return nil
}
