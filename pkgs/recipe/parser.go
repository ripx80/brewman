package recipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Malt struct {
	Name   string
	Unit   string
	Amount float32
}

// internal struct
type Recipe struct {
	Malt []Malt
}

type Parser interface {
	Load(s string) (Recipe, error)
}

// one internal data struct. not the external json!
// interface parse and convert!

func (RecipeM3) Load(s string) (*RecipeM3, error) {
	recipe := &RecipeM3{}
	//rec := &Recipe{}

	// err := json.Unmarshal([]byte(s), recipe)
	// if err != nil {
	// 	return nil, err
	// }
	// recipe.original = s
	// fmt.Println(s)

	// do you will implement all this creepy json stuff? NO

	var result map[string]interface{}
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, err
	}

	// need a sorted map, need malz in correct order..
	m3malt := [7]Malt{}
	for k, v := range result {
		// check type is it a string?
		switch k[:4] {
		//Malz1, Malz1_Menge, Malz1_Einheit
		case "Malz":
			fmt.Println(k)

			num := int(k[4])
			if num < 1 && num > 7 {
				return nil, errors.New("Invalid recipe: only 1 - 7 Malts supported in M3 Recipes")
			}
			if len(k) > 4 {
				m3malt[num].Name = v.(string)
			}
			switch k[5:] {
			case "_Menge":

			case "_Einheit":

			}

		default:
			fmt.Printf("Ignore this field: %s\n", k)
		}

		// switch k {
		// no malt
		// }
	}

	fmt.Println(result)
	return recipe, nil
}

func LoadFile(filename string) (*RecipeM3, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	m3 := RecipeM3{}
	recipe, err := m3.Load(string(content))
	//recipe, err := LoadM3(string(content))
	// convert to Recipe interface
	if err != nil {
		return nil, fmt.Errorf("parsing recipe file %s: %v", filename, err)
	}
	return recipe, nil
}

func (r RecipeM3) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}
