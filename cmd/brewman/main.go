package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	logrusBrave "github.com/ripx80/brave/log/logrus"
	"github.com/ripx80/brave/work"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/recipe"
	"github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	validator "gopkg.in/validator.v2"
)

type configCmd struct {
	configFile   string
	outputFormat string
	debug        *bool
	recipe       *os.File
	kettle       string
}

func absolutePath(fp *os.File) (string, error) {
	return filepath.Abs(fp.Name())
}

func confirm(msg string) bool {
	opt := "use y/n"
	for {
		var response string
		log.Info(msg)
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
		case "y":
			fallthrough
		case "yes":
			return true
		case "n":
			fallthrough
		case "no":
			return false
		default:
			log.WithFields(log.Fields{
				"option": opt,
			}).Warn("no option")
		}

	}

}

func goExit(signals chan os.Signal) {
	signals <- syscall.SIGINT // stops all threats and do a cleanup
	select {}
}

/*
Init setup a kettle with information from a config file
*/
func Init(k *brew.Kettle, kettleConfig config.PodConfig) error {

	var err error
	if kettleConfig == (config.PodConfig{}) {
		return fmt.Errorf("no podconfig in config file. you must have a podconfig to mash/hotwater/cooking")
	}

	// Control section
	switch kettleConfig.Control.Device {
	case "dummy":
		k.Heater = &brew.SSRDummy{}
	case "gpio":
		k.Heater, err = brew.SSRReg(kettleConfig.Control.Address)
		if err != nil {
			return err
		}
	case "external":
		_, err = os.Stat(kettleConfig.Control.Address)
		if err != nil {
			return err
		}
		k.Heater = &brew.External{Cmd: kettleConfig.Control.Address}
	default:
		return fmt.Errorf("unsupported control device: %s", kettleConfig.Control.Device)
	}

	// Agiator section
	switch kettleConfig.Agiator.Device {
	case "dummy":
		k.Agitator = &brew.SSRDummy{}
	case "gpio":
		k.Agitator, err = brew.SSRReg(kettleConfig.Control.Address)
		if err != nil {
			return err
		}
	case "external":
		_, err = os.Stat(kettleConfig.Control.Address)
		if err != nil {
			return err
		}
		k.Agitator = &brew.External{Cmd: kettleConfig.Control.Address}
	case "":
		// can be null no error
		if kettleConfig.Agiator.Address != "" {
			return fmt.Errorf("failed setup agiator, device not set: %s", kettleConfig.Agiator.Address)
		}
	default:
		return fmt.Errorf("unsupported agiator device: %s", kettleConfig.Agiator.Device)
	}

	// Temperatur
	switch kettleConfig.Temperatur.Device {
	case "ds18b20":
		k.Temp, err = brew.DS18B20Reg(kettleConfig.Temperatur.Address)
		if err != nil {
			return err
		}
	case "dummy":
		k.Temp = &brew.TempDummy{Name: "tempdummy", Fn: k.Heater.State, Temp: 20}
	case "default":
		return fmt.Errorf("unsupported temp device: %s", kettleConfig.Temperatur.Device)
	}
	return nil
}

/*Validate check all control instances and run a check programm with increase the water temp*/
func validate(stop chan struct{}, config config.PodConfig) error {
	var err error

	kettle := &brew.Kettle{}
	if err = Init(kettle, config); err != nil {
		return err
	}
	tempTo, err := kettle.Temp.Get()
	if err != nil {
		return err
	}
	if err = kettle.TempUp(stop, tempTo+1); err != nil {
		return err
	}
	return nil
}

func main() {
	defer exit.Safe()

	// config for cmd flags
	cfg := configCmd{}

	a := kingpin.New("brewman", "A command-line brew application")
	a.Version("1.0")
	a.HelpFlag.Short('h')
	a.Author("https://github/ripx80")

	a.Flag("config.file", "brewman configuration file path.").
		StringVar(&cfg.configFile)

	a.Flag("output.format", "output format").
		HintOptions("text", "json").
		Default("text").StringVar(&cfg.outputFormat)

	cfg.debug = a.Flag("output.debug", "Enable debug mode.").Short('v').Bool()

	sc := a.Command("get", "get basic output")
	sc.Command("config", "output current config")
	sc.Command("recipe", "output control information")

	// save in config file
	sc = a.Command("set", "set values")
	sc.Command("config", "save current config to file")
	sr := sc.Command("recipe", "set recipe to brew")
	sr.Arg("filename", "file of the recipe").Required().FileVar(&cfg.recipe)

	sc = a.Command("mash", "mash brew steps")
	sc.Command("start", "start the mash precedure")
	sr = sc.Command("rast", "jump to specific rast")
	rastNum := sr.Arg("num", "rast number [1-8]").Required().Int()

	sc = a.Command("hotwater", "make hotwater in kettle")
	sc.Command("start", "start the hotwater precedure")

	sc = a.Command("cook", "cooking your stuff")
	sc.Command("start", "start the cooking precedure")

	sc = a.Command("control", "control hardware")
	sr = sc.Command("off", "stop all actions")
	sr.Arg("kettle", "stop only actions on kettle (hotwater, masher, cooker)").HintOptions("hotwater", "masher", "cooker").StringVar(&cfg.kettle)
	sr = sc.Command("on", "turn on all")
	sr.Arg("kettle", "turn on kettle (hotwater, masher, cooker)").HintOptions("hotwater", "masher", "cooker").StringVar(&cfg.kettle)
	sr = sc.Command("validate", "validate all devices and run a test program")
	sr.Arg("kettle", "turn on kettle (hotwater, masher, cooker)").HintOptions("hotwater", "masher", "cooker").StringVar(&cfg.kettle)

	logr := logrus.New()

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("parsing commandline arguments")
		a.Usage(os.Args[1:])
	}

	// default config
	configFile, _ := config.Load("")

	if cfg.configFile != "" {
		configFile, err = config.LoadFile(cfg.configFile)
		if err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"configFile": cfg.configFile,
			}).Error("can not load configuration file")
			exit.Exit(1)
		}

	} else {
		cfg.configFile = "brewman.yaml"
		if _, err := os.Stat(cfg.configFile); err == nil {
			configFile, err = config.LoadFile(cfg.configFile)
			if err != nil {
				log.WithFields(log.Fields{
					"error":      err,
					"configFile": cfg.configFile,
				}).Error("can not load configuration file")
				exit.Exit(1)
			}
		}
	}
	// setting up the logger
	if *cfg.debug {
		logr.SetLevel(logrus.DebugLevel)
	}

	if cfg.outputFormat == "json" {
		jf := logrus.JSONFormatter{}
		//jf.PrettyPrint = true
		logr.SetFormatter(&jf)
	}

	if err := validator.Validate(configFile); err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"configFile": cfg.configFile,
		}).Error("config file validation failed")
	}

	// use a wrapper for the WithFields func
	log.Set(logrusBrave.Configured(logr))

	// threads, add data chan, error chan
	signals := make(chan os.Signal, 1)
	stop := make(chan struct{})
	done := make(chan struct{})

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	wg := new(sync.WaitGroup)

	// workaround
	handle := false

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {

	case "set config":
		configFile.Save(cfg.configFile)
		fallthrough

	case "get config":
		fmt.Print(configFile)

	case "set recipe":
		configFile.Recipe.File, err = absolutePath(cfg.recipe)
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"recipe": cfg.recipe,
				"config": cfg.configFile,
			}).Error("can not get recipe file")
			exit.Exit(1)
		}
		// todo parsing recipe and check content

		if err = configFile.Save(cfg.configFile); err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"configFile": cfg.configFile,
			}).Error("set recipe in configuration")
			exit.Exit(1)
		}

	case "get recipe":
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"recipe": configFile.Recipe.File,
			}).Error("get recipe")
			exit.Exit(1)
		}
		fmt.Print(recipe)

		// do this in a function with return err and log in main
	case "control off":
		switch cfg.kettle {
		case "hotwater":
			log.WithFields(log.Fields{
				"kettle": "hotwater",
			}).Info("control off")
			if err := ControlOff(configFile.Hotwater); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "hotwater",
				}).Error("control off failed")
			}
		case "masher":
			log.WithFields(log.Fields{
				"kettle": "masher",
			}).Info("control off")
			if err := ControlOff(configFile.Masher); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "masher",
				}).Error("control off failed")
			}
		case "cooker":
			log.WithFields(log.Fields{
				"kettle": "cooker",
			}).Info("control off")

			if err := ControlOff(configFile.Cooker); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "cooker",
				}).Error("control off failed")
			}
		default:
			log.WithFields(log.Fields{
				"kettle": "masher",
			}).Info("stop all kettle")

			if err := ControlOff(configFile.Hotwater); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "hotwater",
				}).Error("control off failed")
			}
			if err := ControlOff(configFile.Masher); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "masher",
				}).Error("control off failed")
			}
			if err := ControlOff(configFile.Cooker); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "cooker",
				}).Error("control off failed")
			}
		}

	case "control on":
		switch cfg.kettle {
		case "hotwater":
			log.WithFields(log.Fields{
				"kettle": "hotwater",
			}).Info("control on")

			if err := ControlOn(configFile.Hotwater); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "hotwater",
				}).Error("control off failed")

			}
		case "masher":
			log.WithFields(log.Fields{
				"kettle": "masher",
			}).Info("control on")
			if err := ControlOn(configFile.Masher); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "masher",
				}).Error("control off failed")
			}
		case "cooker":
			log.WithFields(log.Fields{
				"kettle": "cooker",
			}).Info("control on")
			if err := ControlOn(configFile.Cooker); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "cooker",
				}).Error("control off failed")
			}
		default:
			log.WithFields(log.Fields{
				"kettle": "hotwater",
			}).Info("turn all kettle on")
			if err := ControlOn(configFile.Hotwater); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "hotwater",
				}).Error("control off failed")
			}
			if err := ControlOn(configFile.Masher); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "masher",
				}).Error("control off failed")
			}

			if err := ControlOn(configFile.Cooker); err != nil {
				log.WithFields(log.Fields{
					"error":  err,
					"kettle": "cooker",
				}).Error("control off failed")
			}
		}

	case "control validate":

		log.Warn("this test not uses recepies. Please test this otherwise")
		switch cfg.kettle {
		case "hotwater":
			if err := validate(stop, configFile.Hotwater); err != nil {
				log.WithFields(log.Fields{
					"kettle": "hotwater",
				}).Error("validation failed")
			}
		case "masher":
			if err := validate(stop, configFile.Masher); err != nil {
				log.WithFields(log.Fields{
					"kettle": "masher",
				}).Error("validation failed")
			}

		case "cooker":
			if err := validate(stop, configFile.Cooker); err != nil {
				log.WithFields(log.Fields{
					"kettle": "cooker",
				}).Error("validation failed")
			}
		default:
			if err := validate(stop, configFile.Hotwater); err != nil {
				log.WithFields(log.Fields{
					"kettle": "hotwater",
				}).Error("validation failed")
			}
			if err := validate(stop, configFile.Masher); err != nil {
				log.WithFields(log.Fields{
					"kettle": "masher",
				}).Error("validation failed")
			}
			if err := validate(stop, configFile.Cooker); err != nil {
				log.WithFields(log.Fields{
					"kettle": "cooker",
				}).Error("validation failed")
			}
		}

	case "hotwater start":
		go func() {
			defer wg.Done()
			wg.Add(1)
			if err := Hotwater(configFile, stop); err != nil {
				log.WithFields(log.Fields{
					"kettle": "hotwater",
					"error":  err,
				}).Error("kettle func failed")
			}
			done <- struct{}{}
		}()

		handle = true

	case "mash start":
		go func() {
			defer wg.Done()
			wg.Add(1)
			if err := Mash(configFile, stop); err != nil {
				log.WithFields(log.Fields{
					"kettle": "masher",
					"error":  err,
				}).Error("kettle func failed")
			}
			done <- struct{}{}
		}()
		handle = true

	case "mash rast":
		go func() {
			defer wg.Done()
			wg.Add(1)
			if err := MashRast(configFile, stop, *(rastNum)); err != nil {
				log.WithFields(log.Fields{
					"kettle": "masher",
					"error":  err,
				}).Error("kettle func failed")
			}
			done <- struct{}{}
		}()

		handle = true

	case "cook start":
		go func() {
			defer wg.Done()
			wg.Add(1)
			if err := Cook(configFile, stop); err != nil {
				log.WithFields(log.Fields{
					"kettle": "cooker",
					"error":  err,
				}).Error("kettle func failed")
			}
			done <- struct{}{}
		}()
		handle = true
	}

	// main handle signals and routines, handle is a workaround
	if handle {
		for {
			select {
			case <-signals:
			case <-stop:
			case <-done:
			}
			close(stop)
			work.WaitTimeout(wg, 1*time.Second) // wait for all workers with timeout
			exit.Exit(0)                        // check the return
		}
	}

}
