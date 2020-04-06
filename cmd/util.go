package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	logrusBrave "github.com/ripx80/brave/log/logrus"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/pod"
	"github.com/ripx80/recipe"
	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

// func absolutePath(fp *os.File) (string, error) {
// 	return filepath.Abs(fp.Name())
// }

func getControl(device, address string, required bool) (brew.Control, error) {
	var (
		control brew.Control
		err     error
	)
	switch device {
	case "dummy":
		control = &brew.SSRDummy{}
	case "gpio":
		control, err = brew.SSRReg(address)
		if err != nil {
			return nil, err
		}
	case "signal":
		code, err := strconv.Atoi(address)
		if err != nil {
			return nil, err
		}
		control = &brew.Signal{Pin: 17, Code: uint64(code)}

	case "external":
		_, err = os.Stat(address)
		if err != nil {
			return nil, err
		}
		control = &brew.External{Cmd: address}
	case "":
		// can be null no error
		if !required && address != "" {
			return nil, fmt.Errorf("failed setup agiator, device not set: %s", address)
		}
		control = &brew.SSRDummy{}

	default:
		return nil, fmt.Errorf("unsupported control device: %s", device)
	}
	return control, nil
}

func getTempSensor(device, address string, state func() bool) (brew.TempSensor, error) {
	var (
		sensor brew.TempSensor
		err    error
	)
	switch device {
	case "ds18b20":
		sensor, err = brew.DS18B20Reg(address)
		if err != nil {
			return nil, err
		}
	case "dummy":
		sensor = &brew.TempDummy{Name: "tempdummy", Fn: state, Temp: 20}
	case "default":
		return nil, fmt.Errorf("unsupported temp device: %s", device)
	}
	return sensor, nil
}

func getKettle(kconf config.PodConfig) (*brew.Kettle, error) {
	var err error
	k := &brew.Kettle{}
	if kconf == (config.PodConfig{}) {
		return nil, fmt.Errorf("no cod config in config file. you must have a cod config to mash/hotwater/cooking")
	}

	// Control Unit
	k.Heater, err = getControl(kconf.Control.Device, kconf.Control.Address, true)
	if err != nil {
		return nil, err
	}

	// Control Unit
	k.Agitator, err = getControl(kconf.Agiator.Device, kconf.Agiator.Address, false)
	if err != nil {
		return nil, err
	}

	// Temperatur
	k.Temp, err = getTempSensor(kconf.Temperatur.Device, kconf.Temperatur.Address, k.Heater.State)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func confirmConsole() error {
	var quest pod.Quest
	opt := "use y/n"
	for {
		select {
		case <-cfg.stop:
			return nil
		case quest = <-cfg.confirm:
			var response string
			fmt.Printf("%s (Y/n)", quest.Msg)
			l, err := fmt.Scan(&response)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("canot read response")
			}
			if l > 3 {
				log.Info(opt) // was a warning
				continue
			}
			response = strings.ToLower(response)
			switch response {
			case "n":
				cfg.confirm <- pod.Quest{Msg: "n", Asw: false}
			default:
				cfg.confirm <- pod.Quest{Msg: "y", Asw: true}
			}

		}
	}
}

func goExit(signals chan os.Signal) {
	signals <- syscall.SIGINT // stops all threats and do a cleanup
	select {}
}

func initChan() {
	// threads, add data chan, error chan
	cfg.signals = make(chan os.Signal, 1)
	cfg.stop = make(chan struct{})
	cfg.confirm = make(chan pod.Quest)
	cfg.done = make(chan struct{}, 2) // buff because after closing stop nobody will recive
	signal.Notify(cfg.signals, syscall.SIGINT, syscall.SIGTERM)
	cfg.wg = new(sync.WaitGroup)
}

func getLogrus() *logrus.Logger {
	logr := logrus.New()
	if cfg.debug {
		logr.SetLevel(logrus.DebugLevel)
	}

	if cfg.format == "json" {
		jf := logrus.JSONFormatter{}
		//jf.PrettyPrint = true
		logr.SetFormatter(&jf)
	}
	return logr
}

func initLogger() {
	// use a wrapper for the WithFields func
	log.Set(logrusBrave.Configured(getLogrus()))
}

func initPods() {
	/*set the recipe per pod: must be set in each cmd and ui per pod*/
	var err error

	hotwater, err := getKettle(cfg.conf.Hotwater)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"kettle": "hotwater",
		}).Error("init failed")
		exit.Exit(1)
	}
	cfg.pods.hotwater = pod.New(hotwater, cfg.recipe)
	cfg.pods.hotwater.Hotwater(cfg.conf.Global.HotwaterTemperatur)

	masher, err := getKettle(cfg.conf.Masher)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"kettle": "masher",
		}).Error("init failed")
		exit.Exit(1)
	}
	cfg.pods.masher = pod.New(masher, cfg.recipe)
	cfg.pods.masher.Mash(cfg.conf.Global.HoldTemperatur, cfg.confirm)

	cooker, err := getKettle(cfg.conf.Cooker)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"kettle": "cooker",
		}).Error("init failed")
		exit.Exit(1)
	}
	cfg.pods.cooker = pod.New(cooker, cfg.recipe)
	cfg.pods.cooker.Cook(cfg.conf.Global.CookingTemperatur)
}

func initRecipe() {
	var err error
	cfg.recipe, err = recipe.LoadFile(cfg.conf.Recipe.File, &recipe.Recipe{})
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"recipe": cfg.conf.Recipe.File,
		}).Error("init recipe")
		exit.Exit(1)
	}
}

func initConfig() {
	cfg.conf, _ = config.Load("")
	var err error

	if cfg.file != "" {
		cfg.conf, err = config.LoadFile(cfg.file)
		if err != nil {
			fmt.Printf("can not load configuration file: %s\n", cfg.file)
			exit.Exit(1)
		}
	} else {
		cfg.file = configfile
		if _, err := os.Stat(cfg.file); err == nil {
			cfg.conf, err = config.LoadFile(cfg.file)
			if err != nil {
				fmt.Printf("can not load configuration file: %s\n", cfg.file)
				exit.Exit(1)
			}
		}
	}
	if err := validator.Validate(cfg.conf); err != nil {
		fmt.Printf("config file validation failed: %s\n", err)
		exit.Exit(1)
	}
}
