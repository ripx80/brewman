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

	sc = a.Command("hotwater", "make hotwater in kettle")
	sc.Command("start", "start the hotwater precedure")

	sc = a.Command("cook", "cooking your stuff")
	sc.Command("start", "start the cooking precedure")

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
		kettle := &brew.Kettle{}
		if err = kettle.Init(configFile.Hotwater); err != nil {
			log.Fatal("Failed to init Kettle:", err)
		}
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("using recipe: ", recipe.Global.Name)
		log.Infof("main water: %f -->  grouting: %f", recipe.Water.MainCast, recipe.Water.Grouting)

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}

		if err := kettle.TempIncreaseTo(configFile.Global.HotwaterTemperatur); err != nil {
			log.Fatal(err)
		}

		err = kettle.TempHolder(configFile.Global.HotwaterTemperatur, 0)
		if err != nil {
			log.Fatal(err)
		}
		// do this in a cleanup func *kettle.cleanup()
		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.Off()
		}

	case "mash start":

		kettle := &brew.Kettle{}
		if err = kettle.Init(configFile.Masher); err != nil {
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

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}
		// Call it here because we must wait for input.. this can take a while
		if err := kettle.TempIncreaseTo(recipe.Mash.InTemperatur); err != nil {
			log.Fatal(err)
		}

		log.Print("Was the malt added? continue: <enter>")
		if !confirm("malt added? continue? <y/n>") {
			if kettle.Agitator != nil && kettle.Agitator.State() {
				kettle.Agitator.Off()
			}
			log.Fatal("Abort!")
		}

		for num, rast := range recipe.Mash.Rests {
			log.Infof("Rast %d: Time: %d Temperatur:%f\n", num, rast.Time, rast.Temperatur)
			if err := kettle.TempIncreaseTo(configFile.Global.HotwaterTemperatur); err != nil {
				log.Error(err)
			}
			err = kettle.TempHolder(rast.Temperatur, time.Duration(rast.Time*60)*time.Second)
			if err != nil {
				log.Error(err)
			}
		}

		if kettle.Agitator != nil && kettle.Agitator.State() {
			kettle.Agitator.Off()
		}
		log.Info("Mashing finished successful")

	case "cook start":
		kettle := &brew.Kettle{}
		if err = kettle.Init(configFile.Cooker); err != nil {
			log.Fatal("Failed to init Kettle:", err)
		}

		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Fatal(err)
		}

		log.Info("using recipe: ", recipe.Global.Name)
		log.Info("cook information: ", recipe.Cook)

		if !confirm("start cooking? <y/n>") {
			os.Exit(0)
		}

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}

		if err := kettle.TempIncreaseTo(configFile.Global.HotwaterTemperatur); err != nil {
			log.Fatal(err)
		}

		err = kettle.TempHolder(configFile.Global.CookingTemperatur, time.Duration(recipe.Cook.Time*60)*time.Second)
		if err != nil {
			log.Fatal(err)

		}

		if kettle.Agitator != nil && kettle.Agitator.State() {
			kettle.Agitator.Off()
		}
		log.Info("Cooking finished successful")

	}
}
