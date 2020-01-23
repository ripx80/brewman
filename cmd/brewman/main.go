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
		log.Infof(msg)
		l, err := fmt.Scan(&response)
		if err != nil {
			log.Errorf("%v", err)
		}
		if l > 3 {
			log.Infof(opt) // was a warning
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
			log.Infof(opt) // was a warn
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
	case "":
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

func validate(stop chan struct{}, config config.PodConfig) error {
	var err error
	// Validate check all control instances and run a check programm with increase the water temp
	kettle := &brew.Kettle{}
	if err = Init(kettle, config); err != nil {
		log.Errorf("Failed to init Kettle: %v", err)
	}
	tempTo, err := kettle.Temp.Get()
	if err != nil {
		log.Errorf("validate: get no temp from sensor")
		return err
	}
	if err = kettle.TempUp(stop, tempTo+1); err != nil {
		log.Errorf("validate: increase temp failed")
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
	//log := log.New(logr)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.Errorf("Error parsing commandline arguments: %v", err)
		a.Usage(os.Args[1:])
	}

	// default config
	configFile, _ := config.Load("")

	if cfg.configFile != "" {
		configFile, err = config.LoadFile(cfg.configFile)
		if err != nil {
			log.Errorf("canot load configuration file: ", err)
			exit.Exit(1)
		}

	} else {
		cfg.configFile = "brewman.yaml"
		if _, err := os.Stat(cfg.configFile); err == nil {
			configFile, err = config.LoadFile(cfg.configFile)
			if err != nil {
				log.Errorf("canot load config file: ", err)
				exit.Exit(1)
			}
		}
	}

	if *cfg.debug {
		logr.SetLevel(logrus.DebugLevel)
	}

	if cfg.outputFormat == "json" {
		jf := logrus.JSONFormatter{}
		//jf.PrettyPrint = true
		logr.SetFormatter(&jf)
	}

	if err := validator.Validate(configFile); err != nil {
		logr.Error("Config file validation failed: ", err)
	}

	log.Set(logr)

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
		log.Infof(fmt.Sprintf("\n%s\n%s", cfg.configFile, configFile))

	case "set recipe":
		configFile.Recipe.File, err = absolutePath(cfg.recipe)
		if err != nil {
			log.Panic(err)

		}
		configFile.Save(cfg.configFile)
		fallthrough

	case "get recipe":
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Panic(err)
		}
		log.Infof(fmt.Sprintf("\n%s\n%s", configFile.Recipe, recipe))

	case "control off":
		switch cfg.kettle {
		case "hotwater":
			log.Infof("stop hotwater")
			if err := ControlOff(configFile.Hotwater); err != nil {
				log.Errorf("Failed to init Kettle:", err)
			}
		case "masher":
			log.Infof("stop masher")
			if err := ControlOff(configFile.Masher); err != nil {
				log.Errorf("Failed to init Kettle:", err)
			}
		case "cooker":
			log.Infof("stop cooker")
			if err := ControlOff(configFile.Cooker); err != nil {
				log.Errorf("Failed to init Kettle:", err)
			}
		default:
			log.Infof("stop all actions and Off")
			if err := ControlOff(configFile.Hotwater); err != nil {
				log.Errorf("Failed to turn off Kettle:", err)
			}
			if err := ControlOff(configFile.Masher); err != nil {
				log.Errorf("Failed to turn off Kettle:", err)
			}
			if err := ControlOff(configFile.Cooker); err != nil {
				log.Errorf("Failed to turn off Kettle:", err)
			}
		}

	case "control on":
		switch cfg.kettle {
		case "hotwater":
			log.Infof("turn on hotwater")
			if err := ControlOn(configFile.Hotwater); err != nil {
				log.Errorf("Failed to turn on Kettle:", err)
			}
		case "masher":
			if err := ControlOn(configFile.Masher); err != nil {
				log.Errorf("Failed to turn on Kettle:", err)
			}
		case "cooker":
			if err := ControlOn(configFile.Cooker); err != nil {
				log.Errorf("Failed to turn on Kettle:", err)
			}
		default:
			log.Infof("turn all kettle on")
			if err := ControlOn(configFile.Hotwater); err != nil {
				log.Errorf("Failed to turn on Kettle:", err)
			}
			if err := ControlOn(configFile.Masher); err != nil {
				log.Errorf("Failed to turn on Kettle:", err)
			}

			if err := ControlOn(configFile.Cooker); err != nil {
				log.Errorf("Failed to turn on Kettle:", err)
			}
		}

	case "control validate":

		log.Infof("this test not uses recepies. Please test this otherwise")
		switch cfg.kettle {
		case "hotwater":
			log.Infof("validate hotwater")
			validate(stop, configFile.Hotwater)
		case "masher":
			log.Infof("validate masher")
			validate(stop, configFile.Masher)
		case "cooker":
			log.Infof("validate cooker")
			validate(stop, configFile.Cooker)
		default:
			log.Infof("validate hotwater")
			validate(stop, configFile.Hotwater)
			log.Infof("validate masher")
			validate(stop, configFile.Masher)
			log.Infof("validate cooker")
			validate(stop, configFile.Cooker)
		}

	case "hotwater start":
		go func() {
			defer wg.Done()
			wg.Add(1)
			if err := Hotwater(configFile, stop); err != nil {
				log.Errorf("error: %v\n", err)
			} else {
				log.Infof("Hotwater finished successful")
			}
			done <- struct{}{}
		}()

		handle = true

	case "mash start":
		go func() {
			defer wg.Done()
			wg.Add(1)
			if err := Mash(configFile, stop); err != nil {
				log.Errorf("error: %v\n", err)
			} else {
				log.Infof("Mashing finished successful")
			}

			done <- struct{}{}
		}()
		handle = true

	case "mash rast":
		go func() {
			defer wg.Done()
			wg.Add(1)
			if err := MashRast(configFile, stop, *(rastNum)); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				log.Infof("Rast finished successful")
			}
			done <- struct{}{}
		}()

		handle = true

	case "cook start":
		go func() {
			defer wg.Done()
			wg.Add(1)
			if err := Cook(configFile, stop); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				log.Infof("Cooking finished successful")
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
