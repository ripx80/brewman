package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     = kingpin.New("brewman", "A command-line brew application")
	output  = app.Flag("o", "output format. use yaml or json").String()
	verbose = app.Flag("v", "verbosity, v=0 quiet, v=1 extended, v=2 debug").Int()

	get        = app.Command("get", "get basic output")
	getConfig  = get.Command("config", "output current config")
	getSensors = get.Command("sensors", "output sensor information")

	set           = app.Command("set", "set values")
	setRecipe     = set.Command("recipe", "set recipe to brew")
	setRecipeFile = setRecipe.Arg("filename", "file of the recipe").Required().File()

	describe = app.Command("describe", "get verbose output of objects")

	validate = app.Command("validate", "validate brewing things like sensors")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "get config":
		fmt.Println("get config")
	case "get sensors":
		fmt.Println("get sensors")
	case "set recipe":
		f := *setRecipeFile
		fmt.Printf("set recipe: %s\n", f.Name())

	}

	fmt.Println(*output)
	fmt.Println(*verbose)

}
