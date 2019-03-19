package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
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

type DefaultConfigCmd struct {
	configFile   string
	outputFormat string
	verbose      int
	recipe       *os.File
}

func main() {

	// config for cmd flags
	cfg := DefaultConfigCmd{}

	//c := config.DefaultConfig

	a := kingpin.New("brewman", "A command-line brew application")
	a.Version("1.0")
	a.HelpFlag.Short('h')
	a.Author("Ripx80")

	a.Flag("config.file", "brewman configuration file path.").
		Default("brewman.yml").StringVar(&cfg.configFile)

	a.Flag("output.format", "output format: yaml, json").
		Default("text").StringVar(&cfg.outputFormat)

	//sc is the tmp placeholder to interact with the subcommand
	sc := a.Command("get", "get basic output")
	sc.Command("config", "output current config")
	sc.Command("sensors", "output sensor information")

	sc = a.Command("set", "set values").Command("recipe", "set recipe to brew")
	sc.Arg("filename", "file of the recipe").Required().FileVar(&cfg.recipe)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {
	case "get config":
		fmt.Println("get config")
	case "get sensors":
		fmt.Println("get sensors")
	case "set recipe":
		fmt.Printf("set recipe: %s\n", cfg.recipe.Name())
		//validate config
		os.Setenv("BREWMAN_RECIPE", "1")
	}

	fmt.Println(cfg.configFile)
	fmt.Println(cfg.outputFormat)
}
