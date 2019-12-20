package recipe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

/*
Recipe as a good json structure for beer recipies
*/
type Recipe struct {
	Global       recipeGlobal       `json:"Global"`
	Water        recipeWater        `json:"Water"`
	Mash         recipeMash         `json:"Mash"`
	Cook         recipeCook         `json:"Cook"`
	Fermentation recipeFermentation `json:"Fermentation"`
	Comments     []Comment          `json:"Comments"`
	original     string
}

type recipeGlobal struct {
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

type recipeWater struct {
	MainCast float64 // Hauptguss
	Grouting float64 // Nachguss in L
}

type recipeMash struct {
	InTemperatur  float64 // Einmaischtemperatur
	OutTemperatur float64 // Abmaischtemperatur
	Malts         []Malt  `json:"Malts"`
	Rests         []Rest  `json:"Rests"`
}

type recipeCook struct {
	Time        int          // in Minutes, Kochzeit
	Ingredients []Ingredient // Weitere Zutaten
	FontHops    []Hop        // Vorderhopfen
	Hops        []Hop
	Whirlpool   []Hop
}

type recipeFermentation struct {
	Yeast       string // Hefe
	Temperatur  float64
	EndDegree   float64      // Endvergärungsgrad
	Carbonation float64      // Karbonisierung
	Hops        []Hop        // Stopfhopfen
	Ingredients []Ingredient // Weitere Zutaten
}

/*
Comment struct to add comments to recepies
*/
type Comment struct {
	Name    string
	Date    string
	Comment string
}

type recipeUnit struct {
	Name   string
	Amount float64
}

type recipeTimeUnit struct {
	recipeUnit
	Time int
}

/*
Malt is a recipe Unit
*/
type Malt = recipeUnit

/*
Ingredient is a recipe TimeUnit
*/
type Ingredient = recipeTimeUnit

/*
Rest holds resting time and temp
*/
type Rest struct {
	Time       int
	Temperatur float64
}

/*
Hop holds unit and Alpha
*/
type Hop struct {
	recipeTimeUnit
	Alpha float64
}

/*
Load unmarshal a recipe
*/
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

/*
PrettyPrint return a pretty string
*/
func (r Recipe) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

/*
Save a recipe to disk
*/
func (r Recipe) Save(fn string) error {
	return ioutil.WriteFile(fn, []byte(r.String()), 0644)
}

/*
SavePretty save a pretty string to disk
*/
func (r Recipe) SavePretty(fn string) error {
	return ioutil.WriteFile(fn, []byte(r.PrettyPrint()), 0644)
}

/*
SavePrettyYaml saves a pretty yamle to disk. not working?
*/
func (r Recipe) SavePrettyYaml(fn string) error {
	s, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fn, []byte(s), 0644)
}

func (r recipeGlobal) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r recipeGlobal) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r recipeMash) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r recipeMash) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r recipeWater) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r recipeWater) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r recipeCook) String() string {
	b, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (r recipeCook) PrettyPrint() string {
	b, err := json.MarshalIndent(r, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}
