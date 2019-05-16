package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/recipe"
	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	validator "gopkg.in/validator.v2"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

type ConfigCmd struct {
	configFile   string
	outputFormat string
	debug        *bool
	recipe       *os.File
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
			log.Error(err)
		}
		if l > 3 {
			log.Warn(opt)
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
			log.Warn(opt)
		}

	}

}

func KettleInit(kettleConfig config.PodConfig) (*brew.Kettle, error) {

	var kettle = &brew.Kettle{}

	if kettleConfig == (config.PodConfig{}) {
		return nil, fmt.Errorf("No Masher in config file. You must have a masher to mash")
	}

	_, err := host.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize periph: %v", err)
	}

	//Heater
	p := gpioreg.ByName(kettleConfig.Control)
	if p == nil {
		return nil, fmt.Errorf("Failed to find Heater Pin: %s", kettleConfig.Agiator)
	}
	kettle.Heater = &brew.SSR{Pin: p}

	// Agiator
	if kettleConfig.Agiator != "" {
		p = gpioreg.ByName(kettleConfig.Agiator)
		if p == nil {
			return nil, fmt.Errorf("Failed to find Agiator Pin: %s", kettleConfig.Agiator)
		}
		kettle.Agitator = &brew.SSR{Pin: p}
	}

	// Temperatur
	kettle.Temp = brew.DS18B20{Name: kettleConfig.Temperatur.Device, Path: kettleConfig.Temperatur.Address}
	if err != nil {
		return nil, fmt.Errorf("Failed to register Temp Sensor: %s", err)
	}
	return kettle, nil
}

func main() {

	// config for cmd flags
	cfg := ConfigCmd{}

	a := kingpin.New("brewman", "A command-line brew application")
	a.Version("1.0")
	a.HelpFlag.Short('h')
	a.Author("Ripx80")

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
	sc.Command("dummy", "dummy the mash precedure")

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.Error("Error parsing commandline arguments: ", err)
		a.Usage(os.Args[1:])
	}

	// default config
	configFile, _ := config.Load("")

	if cfg.configFile != "" {
		configFile, err = config.LoadFile(cfg.configFile)
		if err != nil {
			log.Error("canot load configuration file: ", err)
			os.Exit(1)
		}

	} else {
		cfg.configFile = "brewman.yaml"
		if _, err := os.Stat(cfg.configFile); err == nil {
			configFile, err = config.LoadFile(cfg.configFile)
			if err != nil {
				log.Error("canot load config file: ", err)
				os.Exit(1)
			}
		}
	}

	if *cfg.debug {
		log.SetLevel(log.DebugLevel)
	}

	if cfg.outputFormat == "json" {
		jf := log.JSONFormatter{}
		//jf.PrettyPrint = true
		log.SetFormatter(&jf)
	}

	if err := validator.Validate(configFile); err != nil {
		log.Error("Config file validation failed: ", err)
	}

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {

	case "set config":
		configFile.Save(cfg.configFile)
		fallthrough

	case "get config":
		log.Info(fmt.Sprintf("\n%s\n%s", cfg.configFile, configFile))

	case "set recipe":
		configFile.Recipe.File, err = absolutePath(cfg.recipe)
		if err != nil {
			log.Error("set recipe error: ", err)
		}
		configFile.Save(cfg.configFile)
		fallthrough

	case "get recipe":
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		log.Info(fmt.Sprintf("\n%s\n%s", configFile.Recipe, recipe))

	case "hotwater start":
		kettle, err := KettleInit(configFile.Hotwater)
		if err != nil {
			log.Fatal("Failed to init Kettle:", err)
		}
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("using recipe: ", recipe.Global.Name)
		log.Info("mash information: ", recipe.Water)

		if !confirm("start heating hotwater? <y/n>") {
			os.Exit(0)
		}

		if kettle.Agitator != nil {
			kettle.Agitator.On()
		}
		// InTemperatur
		if err := kettle.ToTemp(configFile.Global.HotwaterTemperatur); err != nil {
			log.Fatal(err)
		}

		// run forever: use os.signals to stop the timer chan,
		// as parameter chan to write to. And build the Timer on main programm
		if err := kettle.HoldTemp(configFile.Global.HotwaterTemperatur, time.Duration(1000000*60)*time.Second); err != nil {
			log.Fatal(err)
		}

		if kettle.Agitator != nil {
			kettle.Agitator.Off()
		}

	case "mash start":

		// todo, change this to a *Kattle way :-)
		kettle, err := KettleInit(configFile.Masher)
		if err != nil {
			log.Fatal("Failed to init Kettle:", err)
		}
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("using recipe: ", recipe.Global.Name)
		log.Info("mash information: ", recipe.Mash)

		if !confirm("start mashing? <y/n>") {
			os.Exit(0)
		}

		if kettle.Agitator != nil {
			kettle.Agitator.On()
		}
		// InTemperatur
		if err := kettle.ToTemp(recipe.Mash.InTemperatur); err != nil {
			log.Fatal(err)
		}

		if kettle.Agitator != nil {
			kettle.Agitator.Off()
		}

		// Give Malts
		log.Print("Was the malt added? continue: <enter>")
		if !confirm("malt added? continue? <y/n>") {
			os.Exit(0)
		}

		if kettle.Agitator != nil {
			kettle.Agitator.On()
		}
		// Step Rests
		cnt := 1
		for _, rast := range recipe.Mash.Rests {
			log.Infof("Rast %d: Time: %d Temperatur:%f\n", cnt, rast.Time, rast.Temperatur)

			if err := kettle.ToTemp(rast.Temperatur); err != nil {
				log.Fatal(err)
			}

			if err := kettle.HoldTemp(rast.Temperatur, time.Duration(rast.Time*60)*time.Second); err != nil {
				log.Fatal(err)
			}
			cnt++
		}
		if kettle.Agitator != nil {
			kettle.Agitator.Off()
		}
		log.Info("Finish Mash")

	case "mash dummy":

		log.Info("check configfile")

		if configFile.Masher == (config.PodConfig{}) {
			log.Error("No Masher in config file. You must have a masher to mash :-)")
			os.Exit(1)
		}

		log.Info("check recipe")
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		log.Infof("using recipe: %s\n", recipe.Global.Name)
		log.Infof("mash information:\n%s", recipe.Mash)

		masher := &brew.Kettle{
			Temp:     &brew.TempDummy{Name: "TempDummy", Fn: func(x float64) float64 { return x + 4.0 }, Temp: 40.0},
			Heater:   &brew.SSRDummy{},
			Agitator: &brew.SSRDummy{},
		}
		if !confirm("start mashing? <y/n>") {
			os.Exit(0)
		}

		// InTemperatur
		if err := masher.ToTemp(recipe.Mash.InTemperatur); err != nil {
			log.Fatal(err)
		}

		// Give Malts
		if !confirm("malt added? continue? <y/n>") {
			os.Exit(0)
		}

		// Step Rests
		cnt := 1
		for _, rast := range recipe.Mash.Rests {
			log.Infof("Rast %d: Time: %d Temperatur:%f\n", cnt, rast.Time, rast.Temperatur)

			masher.Temp = &brew.TempDummy{Name: "TempDummy", Fn: func(x float64) float64 { return x + 4.0 }, Temp: 40.0}
			if err := masher.ToTemp(rast.Temperatur); err != nil {
				log.Fatal(err)
			}
			masher.Temp = &brew.TempDummy{Name: "TempDummy", Fn: brew.UpDown, Temp: 55.0}
			if err := masher.HoldTemp(rast.Temperatur, time.Duration(rast.Time-(rast.Time-4))*time.Second); err != nil {
				log.Fatal(err)
			}
			cnt++
		}
		log.Info("Finish Mash")
	}

}
