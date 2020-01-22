package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	log "github.com/ripx80/brave/log/logger"
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

func validate(stop chan struct{}, kettle *brew.Kettle, config config.PodConfig) error {
	var err error
	// Validate check all control instances and run a check programm with increase the water temp
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
	log := log.New(logr)

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
			os.Exit(1)
		}

	} else {
		cfg.configFile = "brewman.yaml"
		if _, err := os.Stat(cfg.configFile); err == nil {
			configFile, err = config.LoadFile(cfg.configFile)
			if err != nil {
				log.Errorf("canot load config file: ", err)
				os.Exit(1)
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

	// threads, add data chan, error chan
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	stop := make(chan struct{})
	kettle := &brew.Kettle{} //need for cleanup
	//wg := new(sync.WaitGroup)

	// signal handler // only for threat jobs
	go func() {
		defer os.Exit(0) //add error channel
		select {
		case <-signals:
		case <-stop:
		}
		close(stop)
		log.Infof("cleanup in controller threat")
		kettle.Off()
		//wg.Wait() // wait for all workers
		log.Infof("go exit")
	}()

	// watch threat, must have temp initilized... must call after a successful init of a Kettle
	// think about a good implementation
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	if err := kettle.Watch(stop, 2); err != nil {
	// 		log.Error(err)
	// 		goExit(signals)
	// 	}
	// }()

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {

	case "set config":
		configFile.Save(cfg.configFile)
		fallthrough

	case "get config":
		log.Infof(fmt.Sprintf("\n%s\n%s", cfg.configFile, configFile))

	case "set recipe":
		configFile.Recipe.File, err = absolutePath(cfg.recipe)
		if err != nil {
			log.Errorf("set recipe error: ", err)
		}
		configFile.Save(cfg.configFile)
		fallthrough

	case "get recipe":
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Errorf("%v", err)
			os.Exit(1)
		}
		log.Infof(fmt.Sprintf("\n%s\n%s", configFile.Recipe, recipe))

	case "control off":
		switch cfg.kettle {
		case "hotwater":
			log.Infof("stop hotwater")
			if err = Init(kettle, configFile.Hotwater); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
		case "masher":
			log.Infof("stop masher")
			if err = Init(kettle, configFile.Masher); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
		case "cooker":
			log.Infof("stop cooker")
			if err = Init(kettle, configFile.Cooker); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
		default:
			log.Infof("stop all actions and Off")
			if err = Init(kettle, configFile.Hotwater); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
			kettle.Off()

			if err = Init(kettle, configFile.Masher); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
			kettle.Off()

			if err = Init(kettle, configFile.Cooker); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
		}
		goExit(signals)
	case "control on":
		switch cfg.kettle {
		case "hotwater":
			log.Infof("turn on hotwater")
			if err = Init(kettle, configFile.Hotwater); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
		case "masher":
			log.Infof("turn on masher")
			if err = Init(kettle, configFile.Masher); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
		case "cooker":
			log.Infof("turn on cooker")
			if err = Init(kettle, configFile.Cooker); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
		default:
			log.Infof("turn all kettle on")
			if err = Init(kettle, configFile.Hotwater); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
			kettle.On()

			if err = Init(kettle, configFile.Masher); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
			kettle.On()

			if err = Init(kettle, configFile.Cooker); err != nil {
				log.Errorf("Failed to init Kettle:", err)
				panic("USE SAFE FUNC")
			}
		}
		kettle.On()
	case "control validate":

		log.Infof("this test not uses recepies. Please test this otherwise")
		switch cfg.kettle {
		case "hotwater":
			log.Infof("validate hotwater")
			validate(stop, kettle, configFile.Hotwater)
		case "masher":
			log.Infof("validate masher")
			validate(stop, kettle, configFile.Masher)
		case "cooker":
			log.Infof("validate cooker")
			validate(stop, kettle, configFile.Cooker)
		default:
			log.Infof("validate hotwater")
			validate(stop, kettle, configFile.Hotwater)
			log.Infof("validate masher")
			validate(stop, kettle, configFile.Masher)
			log.Infof("validate cooker")
			validate(stop, kettle, configFile.Cooker)
		}

		// threat and temp watcher work
	case "hotwater start":

		if err = Init(kettle, configFile.Hotwater); err != nil {
			log.Errorf("Failed to init Kettle:", err)
			panic("USE SAFE FUNC")
		}
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Errorf("%v", err)
			os.Exit(1)
		}

		log.Infof("using recipe: ", recipe.Global.Name)
		log.Infof("main water: %f -->  grouting: %f", recipe.Water.MainCast, recipe.Water.Grouting)

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}

		if err := kettle.TempUp(stop, configFile.Global.HotwaterTemperatur); err != nil {
			log.Errorf("%v", err)
			goExit(signals)
		}

		if err := kettle.TempHold(stop, configFile.Global.HotwaterTemperatur, 0); err != nil {
			log.Errorf("%v", err)
			goExit(signals)
		}

		goExit(signals)

	case "mash start":
		if err = Init(kettle, configFile.Masher); err != nil {
			log.Errorf("Failed to init Kettle:", err)
			panic("USE SAFE FUNC")
		}

		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Errorf("%v", err)
			panic("USE SAFE FUNC")
			//os.Exit(1)
		}

		log.Infof("using recipe: ", recipe.Global.Name)
		log.Infof("mash information: ", recipe.Mash)

		if !confirm("start mashing? <y/n>") {
			goExit(signals)
		}

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}

		if err := kettle.TempUp(stop, recipe.Mash.InTemperatur); err != nil {
			log.Errorf("%v", err)
			goExit(signals)
		}

		if !confirm("malt added? continue? <y/n>") {
			goExit(signals)
		}

		for num, rast := range recipe.Mash.Rests {
			log.Infof("Rast %d: Time: %d Temperatur:%f\n", num, rast.Time, rast.Temperatur)

			if err := kettle.TempUp(stop, rast.Temperatur); err != nil {
				log.Errorf("%v", err)
				goExit(signals)
			}

			if err := kettle.TempHold(stop, rast.Temperatur, time.Duration(rast.Time*60)*time.Second); err != nil {
				log.Errorf("%v", err)
				goExit(signals)
			}

		}

		log.Infof("Mashing finished successful")
		goExit(signals)

	case "mash rast":

		num := *(rastNum)

		if err = Init(kettle, configFile.Masher); err != nil {
			log.Errorf("Failed to init Kettle:", err)
			panic("USE SAFE FUNC")
		}

		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Errorf("%v", err)
			panic("USE SAFE FUNC")
			//os.Exit(1)
		}

		if num > 8 || num <= 0 {
			log.Errorf("rast number out of range [1-8]")
			panic("USE SAFE FUNC")
			//os.Exit(1)
		}

		if len(recipe.Mash.Rests) < num {
			log.Errorf("rast number not in recipe")
			panic("USE SAFE FUNC")
			//os.Exit(1)
		}

		log.Infof("jump to rast number: %d", num)
		log.Infof("using recipe: ", recipe.Global.Name)

		rast := recipe.Mash.Rests[num-1]
		log.Infof("Rast %d: Time: %d Temperatur: %.2f\n", num, rast.Time, rast.Temperatur)
		if err := kettle.TempUp(stop, rast.Temperatur); err != nil {
			log.Errorf("%v", err)
			goExit(signals)
		}

		if err := kettle.TempHold(stop, rast.Temperatur, time.Duration(rast.Time*60)*time.Second); err != nil {
			log.Errorf("%v", err)
			goExit(signals)
		}

		log.Infof("Rast finished successful")
		goExit(signals)

	case "cook start":
		if err = Init(kettle, configFile.Cooker); err != nil {
			log.Errorf("Failed to init Kettle:", err)
			panic("USE SAFE FUNC")
		}

		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Errorf("%v", err)
			panic("USE SAFE FUNC")
		}

		log.Infof("using recipe: ", recipe.Global.Name)
		log.Infof("cook information: ", recipe.Cook)

		if !confirm("start cooking? <y/n>") {
			goExit(signals)
		}

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}

		if err := kettle.TempUp(stop, configFile.Global.CookingTemperatur); err != nil {
			log.Errorf("%v", err)
			goExit(signals)
		}

		if err := kettle.TempHold(stop, configFile.Global.CookingTemperatur, time.Duration(recipe.Cook.Time*60)*time.Second); err != nil {
			log.Errorf("%v", err)
			goExit(signals)
		}

		log.Infof("Cooking finished successful")
		goExit(signals)
	}
}
