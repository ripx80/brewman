package main

import (
	"fmt"
	"os"

	"github.com/ripx80/brewman/pkgs/recipe"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	a := kingpin.New("m3c", "A command-line converter for brew recepies maischemalzundmehr")
	a.Version("1.0")
	a.HelpFlag.Short('h')
	a.Author("Ripx80")

	format := a.Flag("o", "output format, yaml/json").
		HintOptions("yaml", "json").
		Default("json").String()

	in := a.Arg("input-file", "Recipe file to convert").Required().String()
	out := a.Arg("output-file", "Recipe file for converted output").Required().String()

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing commandline arguments: ", err)
		a.Usage(os.Args[1:])
	}

	recipe, err := recipe.LoadFile(*in, &recipe.RecipeM3{})
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	fmt.Println(recipe)
	if *format == "yaml" {
		err = recipe.SavePrettyYaml(*out)

	} else {
		err = recipe.SavePretty(*out)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
