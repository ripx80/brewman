package recipe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type Recipe struct {
	Global       RecipeGlobal       `json:"Global"`
	Water        RecipeWater        `json:"Water"`
	Mash         RecipeMash         `json:"Mash"`
	Cook         RecipeCook         `json:"Cook"`
	Fermentation RecipeFermentation `json:"Fermentation"`
	Comment      []RecipeComment    `json:"Comments"`
	original     string
}

type RecipeGlobal struct {
	Name              string  `json:"Name" validate:"nonzero"`
	Date              string  `json:"Date" validate:"nonzero"`
	Sort              string  `json:"Sort" validate:"nonzero"`
	Author            string  `json:"Author" validate:"nonzero"`
	Clone             string  `json:"Clone" validate:"nonzero"`
	DecisiveSeasoning float64 // Ausschlagwürze
	SudYield          float64 // Sudausbeute
	OriginalWort      float64 //Stammwürze
	IBU               int
	Color             int
	Alcohol           float64
	ShortDescription  string
	Annotation        string // Anmerkungen des Authors
}

type RecipeWater struct {
	MainCast float64 // Hauptguss
	Grouting float64 // Nachguss in L
}

type RecipeMash struct {
	InTemperatur  float64 // Einmaischtemperatur
	OutTemperatur float64 // Abmaischtemperatur
	Malts         []Malt  `json:"Malts"`
	Rests         []Rest  `json:"Rests"`
}

type RecipeCook struct {
	Time        int          // in Minutes, Kochzeit
	Ingredients []Ingredient // Weitere Zutaten
	FontHops    []Hop        // Vorderhopfen
	Hops        []Hop
	Whirpool    []Hop
}

type RecipeFermentation struct {
	Yeast       string // Hefe
	Temperatur  float64
	EndDegree   float64      // Endvergärungsgrad
	Carbonation float64      // Karbonisierung
	Hops        []Hop        // Stopfhopfen
	Ingredients []Ingredient // Weitere Zutaten
}

type RecipeComment struct {
	Name    string
	Date    string
	Comment string
}

type RecipeUnit struct {
	Name   string
	Amount float64
}

type RecipeTimeUnit struct {
	RecipeUnit
	Time int
}

type Malt = RecipeUnit
type Ingredient = RecipeTimeUnit

type Rest struct {
	Time       int
	Temperatur float64
}

type Hop struct {
	RecipeTimeUnit
	Alpha float64
}

func (r *Recipe) Load(s string) (*Recipe, error) {
	err := yaml.Unmarshal([]byte(s), r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r Recipe) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r Recipe) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r Recipe) Save(fn string) error {
	return ioutil.WriteFile(fn, []byte(r.String()), 0644)
}

func (r Recipe) SavePretty(fn string) error {
	return ioutil.WriteFile(fn, []byte(r.PrettyPrint()), 0644)
}

//not working
func (r Recipe) SavePrettyYaml(fn string) error {
	s, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fn, []byte(s), 0644)
}

func (r RecipeGlobal) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r RecipeGlobal) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r RecipeMash) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r RecipeMash) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r RecipeWater) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r RecipeWater) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r RecipeCook) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r RecipeCook) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}
