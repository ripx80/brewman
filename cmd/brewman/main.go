package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/recipe"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/validator.v2"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type ConfigCmd struct {
	configFile   string
	outputFormat string
	verbose      int
	recipe       *os.File
}

func absolutePath(fp *os.File) (string, error) {
	return filepath.Abs(fp.Name())
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

	//sc is the tmp placeholder to interact with the subcommand
	sc := a.Command("get", "get basic output")
	sc.Command("config", "output current config")
	//sc.Command("sensors", "output sensor information")
	//sc.Command("control", "output control information")
	sc.Command("recipe", "output control information")

	// save in config file
	sc = a.Command("set", "set values")
	sc.Command("config", "save current config to file")
	sr := sc.Command("recipe", "set recipe to brew")
	sr.Arg("filename", "file of the recipe").Required().FileVar(&cfg.recipe)

	sc = a.Command("start", "start brew steps")
	sc.Command("mash", "start the mash precedure")

	//add sensor?
	//delete sensor?

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

	if cfg.outputFormat == "json" {
		jf := log.JSONFormatter{}
		//jf.PrettyPrint = true
		log.SetFormatter(&jf)
	}

	log.SetLevel(log.InfoLevel)

	if err := validator.Validate(configFile); err != nil {
		log.Error("Config file validation failed: ", err)
	}

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {

	case "set config":
		configFile.Save(cfg.configFile)
		fallthrough

	case "get config":
		log.Info(fmt.Sprintf("\n%s\n%s", cfg.configFile, configFile))

	// case "get sensors":
	// 	log.Info(configFile.Sensor)

	// case "get controls":
	// 	log.Info(configFile.Control)

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

	case "start mash":

		// check if masher is configured
		if configFile.Masher == (config.PodConfig{}) {
			log.Error("No Masher in config file. You must have a masher to mash :-)")
			os.Exit(1)
		}

		// brew.Init() // init all devices and sensors aso
		per := &brew.Periph{TempSensors: make(map[string]brew.TempSensor), Controls: make(map[string]brew.Control)}
		err := per.Init()
		if err != nil {
			log.Error(err)
		}

		/*
			HotTube.init(TempSensor, Control)
			Masher.Init(TempSensor, Control)
			Cooker.Init(TempSensor, Control)

			TempSensor: Name, Bus, Address
		*/

		// use periph/cmd/onewire-list to get all informations

		// init the masher
		ssr := &brew.SSR{Pin: gpioreg.ByName(configFile.Masher.Control)}
		per.Controls["Masher-Control"] = ssr

		ssr = &brew.SSR{Pin: gpioreg.ByName(configFile.Masher.Agiator)}
		per.Controls["Masher-Agitator"] = ssr

		//ds.Device, err = ds18b20.New(&bus, addr, 10)
		ds := &brew.DS18B20{}
		err = ds.InitDummy()
		if err != nil {
			log.Error("Failed to register Temp Sensor")
			os.Exit(1)
		}
		per.TempSensors["Masher-Temperatur"] = ds

		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		masher := &brew.Masher{
			Temp:     per.TempSensors["Masher-Temperatur"],
			Heater:   per.Controls["SSR-Plate-1"],
			Agitator: per.Controls["SSR-Agitator"],
			Recipe:   recipe.Mash,
		}

		logc := make(chan string)
		done := make(chan error, 2)
		defer close(done)

		go func() {
			done <- masher.Mash(logc)
		}()

		go func() {
			done <- func() error {
				for {
					j, more := <-logc
					if more {
						log.Info(j)
					} else {
						return nil
					}
				}
			}()
		}()

		var stopped bool
		for i := 0; i < cap(done); i++ {
			if err := <-done; err != nil {
				log.Error(err)
			}
			if !stopped {
				stopped = true
				close(logc)
				log.Info("Close")
			}
		}
		log.Info("Q")
		//https://github.com/heptio/workgroup

	}

}
