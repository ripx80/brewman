package main

import (
	"os"
	"path/filepath"

	"github.com/ripx80/brewman/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// var (
// 	app = kingpin.New(os.Args[0], "A command-line brew application")

// 	output  = app.Flag("o", "output format. use yaml or json").String()
// 	verbose = app.Flag("v", "verbosity, v=0 quiet, v=1 extended, v=2 debug").Int()

// 	get        = app.Command("get", "get basic output")
// 	getConfig  = get.Command("config", "output current config")
// 	getSensors = get.Command("sensors", "output sensor information")

// 	set           = app.Command("set", "set values")
// 	setRecipe     = set.Command("recipe", "set recipe to brew")
// 	setRecipeFile = setRecipe.Arg("filename", "file of the recipe").Required().File()

// 	describe = app.Command("describe", "get verbose output of objects")

// 	validate = app.Command("validate", "validate brewing things like sensors")
// )

type ConfigCmd struct {
	configFile   *os.File
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
		FileVar(&cfg.configFile)

	a.Flag("output.format", "output format").
		HintOptions("text", "json").
		Default("text").StringVar(&cfg.outputFormat)

	//sc is the tmp placeholder to interact with the subcommand
	sc := a.Command("get", "get basic output")
	sc.Command("config", "output current config")
	sc.Command("sensors", "output sensor information")
	sc.Command("control", "output control information")
	sc.Command("recipe", "output control information")

	// save in config file
	sc = a.Command("set", "set values").Command("recipe", "set recipe to brew")
	sc.Arg("filename", "file of the recipe").Required().FileVar(&cfg.recipe)

	//add sensor?
	//delete sensor?

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.Error("Error parsing commandline arguments: ", err)
		a.Usage(os.Args[1:])
	}

	// default config if no config file is present
	configFile, err := config.Load("")

	if cfg.configFile == nil {
		fp, err := os.Open("brewman.yaml")
		if err == nil {
			cfg.configFile = fp
		}
	}

	if cfg.configFile != nil {
		fp, err := filepath.Abs(cfg.configFile.Name())
		if err != nil {
			log.Error("Error parsing config file path: ", err)
		}
		configFile, err = config.LoadFile(fp)
	}

	if cfg.outputFormat == "json" {
		jf := log.JSONFormatter{}
		//jf.PrettyPrint = true
		log.SetFormatter(&jf)
	}

	log.SetLevel(log.InfoLevel)

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {
	case "get config":
		log.Info(configFile)

	case "get sensors":
		log.Info(configFile.Sensor)

	case "get controls":
		log.Info(configFile.Control)

	case "set recipe":

		configFile.Recipe.File, err = absolutePath(cfg.recipe)
		if err != nil {
			log.Error("set recipe error: ", err)
		}
		//todo: parse and validate, append to config file

		fallthrough
	case "get recipe":
		log.Info(configFile.Recipe)

	}

}
