package recipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

// internal struct

type Parser interface {
	Load(s string) (Recipe, error)
	Convert() (Recipe, error)
}

func (rm *RecipeM3) UnmarshalJSON(data []byte) error {

	// Unmarshal good data from this creepy json.
	type plain RecipeM3
	if err := json.Unmarshal(data, ((*plain)(rm))); err != nil {
		return err
	}

	// next all data will be parsed manually
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	// Malt
	cnt := 1
	for KeyExists(result, fmt.Sprintf("Malz%d", cnt)) {
		k := fmt.Sprintf("Malz%d", cnt)
		malt := &Malt{}

		switch result[k].(type) {
		case string:
			malt.Name = result[k].(string)
		default:
			return errors.New("incorrect json type: malt name must be a string")
		}

		k = fmt.Sprintf("%s_Menge", k)
		switch result[k].(type) {
		case float32, float64:
			malt.Amount = result[k].(float64)
		default:
			return errors.New("incorrect json type: malt amount must be type of float")
		}

		rm.Malts = append(rm.Malts, *malt)
		cnt++
	}

	/*
		Rests
		"Infusion_Rasttemperatur1": "52",
		"Infusion_Rastzeit1": "15",
	*/
	cnt = 1
	for KeyExists(result, fmt.Sprintf("Infusion_Rasttemperatur%d", cnt)) {
		k := fmt.Sprintf("Infusion_Rasttemperatur%d", cnt)
		rest := &Rest{}

		switch result[k].(type) {
		case string:
			v := result[k].(string)
			if v == "" {
				cnt++
				continue
			}
			if rest.Temperatur, err = strconv.ParseFloat(v, 64); err != nil {
				return err
			}
		default:
			return errors.New("incorrect json type: malt name must be a string")
		}

		k = fmt.Sprintf("Infusion_Rastzeit%d", cnt)
		if !KeyExists(result, k) {
			return errors.New("incorrect json: found Rasttemp but no Rastzeit in recipe")
		}

		switch result[k].(type) {
		case string:
			if rest.Time, err = strconv.Atoi(result[k].(string)); err != nil {
				return err
			}
		default:
			return errors.New("incorrect json type: rast temp must be a string")
		}
		rm.Rests = append(rm.Rests, *rest)
		cnt++
	}

	// here I was lazy....
	/*Font Hops
	"Hopfen_VWH_1_Sorte": "",
	"Hopfen_VWH_1_Menge": "",
	"Hopfen_VWH_1_alpha": "",
	*/
	cnt = 1
	for KeyExists(result, fmt.Sprintf("Hopfen_VWH_%d_Sorte", cnt)) {
		k := fmt.Sprintf("Hopfen_VWH_%d_Sorte", cnt)
		hop := &Hop{}

		if hop.Name = result[k].(string); hop.Name == "" {
			cnt++
			continue
		}

		k = fmt.Sprintf("Hopfen_VWH_%d_Menge", cnt)
		v := result[k].(string)
		if !KeyExists(result, k) || v == "" {
			return errors.New("Font Hop found but no amount key or value!")
		}
		if hop.Amount, err = strconv.Atoi(v); err != nil {
			return err
		}

		k = fmt.Sprintf("Hopfen_VWH_%d_alpha", cnt)
		v = result[k].(string)
		if !KeyExists(result, k) || v == "" {
			return errors.New("Font Hop found but no alpha!")
		}
		if hop.Alpha, err = strconv.ParseFloat(result[k].(string), 64); err != nil {
			return err
		}
		rm.FontHops = append(rm.FontHops, *hop)
		cnt++
	}
	// cnt = 1
	// for KeyExists(result, fmt.Sprintf("Hopfen_%d_Sorte", cnt)) {
	// 	k := fmt.Sprintf("Hopfen_%d_Sorte", cnt)
	// 	hop := &Hop{}

	// 	if hop.Name = result[k].(string); hop.Name == "" {
	// 		cnt++
	// 		continue
	// 	}

	/*
				"Hopfen_1_Kochzeit": "70",
				"Hopfen_1_Sorte": "Fuggles, Pellets",
		    	"Hopfen_1_Menge": "280",
				"Hopfen_1_alpha": "4.1",


				"WeitereZutat_Wuerze_1_Name": "",
		    	"WeitereZutat_Wuerze_1_Menge": "",
		    	"WeitereZutat_Wuerze_1_Einheit": "g",
				"WeitereZutat_Wuerze_1_Kochzeit": "",

				"Stopfhopfen_1_Sorte": "Fuggles, Dolden ",
				"Stopfhopfen_1_Menge": "55",

				"WeitereZutat_Gaerung_1_Name": "",
		    	"WeitereZutat_Gaerung_1_Menge": "",
				"WeitereZutat_Gaerung_1_Einheit": "g",

				"Hopfen_3_Kochzeit": "Whirlpool",
	*/

	return nil
}

func KeyExists(m map[string]interface{}, k string) bool {
	if _, ok := m[k]; ok {
		return true
	}
	return false
}

func (RecipeM3) Load(s string) (*RecipeM3, error) {
	recipe := &RecipeM3{}

	err := json.Unmarshal([]byte(s), recipe)
	if err != nil {
		return nil, err
	}
	recipe.original = s
	fmt.Println(s)

	// need a sorted map, need only rast in correct order..

	// for k, v := range result {
	// 	// check type is it a string?
	// 	switch k[:4] {
	// 	//Malz1, Malz1_Menge, Malz1_Einheit
	// 	case "Malz":
	// 		fmt.Println(k)

	// 		num := int(k[4])
	// 		if num < 1 && num > 7 {
	// 			return nil, errors.New("Invalid recipe: only 1 - 7 Malts supported in M3 Recipes")
	// 		}
	// 		if len(k) > 4 {
	// 			m3malt[num].Name = v.(string)
	// 		}
	// 		switch k[5:] {
	// 		case "_Menge":

	// 		case "_Einheit":

	// 		}

	// 	default:
	// 		fmt.Printf("Ignore this field: %s\n", k)
	// 	}

	// switch k {
	// no malt
	// }
	//}

	fmt.Println(recipe)
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
