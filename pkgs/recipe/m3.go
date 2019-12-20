package recipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type converter struct {
	keys map[string]string
	cmap map[string]interface{}
	pos  int
}

type m3Global struct {
	Name   string `json:"Name" validate:"nonzero"`
	Date   string `json:"Datum" validate:"nonzero"`
	Sort   string `json:"Sorte" validate:"nonzero"`
	Author string `json:"Autor" validate:"nonzero"`
	//	Clone string `json:"Klonbier" validate:"nonzero"`

	Clone             string  `json:"Klonbier_Original"`
	DecisiveSeasoning float64 `json:"Ausschlagswuerze,float64"`
	SudYield          float64 `json:"Sudhausausbeute,float64"`
	OriginalWort      float64 `json:"Stammwuerze,float64"`
	IBU               int     `json:"Bittere"`
	Color             int     `json:"Farbe,string"`
	Alcohol           float64 `json:"Alkohol"`
	ShortDescription  string  `json:"Kurzbeschreibung"`
	Annotation        string  `json:"Anmerkung_Autor,string"`
}

/*
M3 implements Maische-Malz-und-Mehr recipies
*/
type M3 struct {
	m3Global

	MainCast      float64 `json:"Infusion_Hauptguss,string"`
	InTemperatur  float64 `json:"Infusion_Einmaischtemperatur"`
	OutTemperatur float64 `json:"Abmaischtemperatur,string"`
	Grouting      float64 `json:"Nachguss,string"`
	Time          int     `json:"Kochzeit_Wuerze,string"`
	Yeast         string  `json:"Hefe"`
	EndDegree     float64 `json:"Endvergaerungsgrad,string"`
	Carbonation   float64 `json:"Karbonisierung,string"`
	Annotation    string  `json:"Anmerkung_Autor"`

	Temperatur   float64
	Malts        []Malt
	Rests        []Rest
	FontHops     []Hop
	Hops         []Hop
	Whirlpool    []Hop
	Ingredients  []Ingredient
	Fermentation recipeFermentation
}

/*
Load read the json file from m3
*/
func (rm *M3) Load(s string) (*Recipe, error) {
	err := json.Unmarshal([]byte(s), rm)
	if err != nil {
		return nil, err
	}
	recipe := &Recipe{}
	//recipe.original = s
	recipe.Global = recipeGlobal{
		Name:              rm.Name,
		Date:              rm.Date,
		Sort:              rm.Sort,
		Author:            rm.Author,
		Clone:             rm.Clone,
		DecisiveSeasoning: rm.DecisiveSeasoning,
		SudYield:          rm.SudYield,
		OriginalWort:      rm.OriginalWort,
		IBU:               rm.IBU,
		Color:             rm.Color,
		Alcohol:           rm.Alcohol,
		ShortDescription:  rm.ShortDescription,
		Annotation:        rm.Annotation,
	}

	recipe.Water = recipeWater{
		MainCast: rm.MainCast,
		Grouting: rm.Grouting,
	}

	recipe.Mash = recipeMash{
		InTemperatur:  rm.InTemperatur,
		OutTemperatur: rm.OutTemperatur,
		Malts:         rm.Malts,
		Rests:         rm.Rests,
	}

	recipe.Cook = recipeCook{
		Time:        rm.Time,
		Ingredients: rm.Ingredients,
		FontHops:    rm.FontHops,
		Hops:        rm.Hops,
		Whirlpool:   rm.Whirlpool,
	}

	recipe.Fermentation = recipeFermentation{
		Yeast:       rm.Yeast,
		Temperatur:  rm.Temperatur,
		EndDegree:   rm.EndDegree,
		Carbonation: rm.Carbonation,
		Hops:        rm.Fermentation.Hops,
		Ingredients: rm.Fermentation.Ingredients,
	}
	return recipe, nil
}

func (rm M3) String() string {
	b, err := json.Marshal(rm)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

/*
PrettyPrint return a pretty string
*/
func (rm M3) PrettyPrint() string {
	b, err := json.MarshalIndent(rm, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

/*
UnmarshalJSON implements the json unmarshaller. The json is difficult to decode. convert json fields
*/
func (rm *M3) UnmarshalJSON(data []byte) error {
	// this is no good implementation. I hope my skills will be better in the future to do this -.-
	// Unmarshal good data from this creepy json.

	type plain M3
	if err := json.Unmarshal(data, ((*plain)(rm))); err != nil {
		return err
	}

	// next all data will be parsed manually
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	conv := &converter{}
	// add a test file for this
	// in this field are different values. normalize...

	re := regexp.MustCompile("[0-9]+")
	x := re.FindAllString(result["Gaertemperatur"].(string), -1)
	var f float64
	for _, value := range x {
		f, err = strconv.ParseFloat(value, 64)
		rm.Temperatur += f
	}
	rm.Temperatur = rm.Temperatur / float64(len(x))
	conv.cmap = result
	conv.keys = map[string]string{"Name": "Malz%d", "Amount": "Malz%d_Menge", "Unit": "Malz%d_Einheit"}
	rm.Malts, err = conv.Malts()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Temperatur": "Infusion_Rasttemperatur%d", "Time": "Infusion_Rastzeit%d"}
	rm.Rests, err = conv.Rests()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Name": "Hopfen_VWH_%d_Sorte", "Amount": "Hopfen_VWH_%d_Menge", "Alpha": "Hopfen_VWH_%d_alpha"}
	rm.FontHops, err = conv.FontHop()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Name": "Hopfen_%d_Sorte", "Amount": "Hopfen_%d_Menge", "Time": "Hopfen_%d_Kochzeit", "Alpha": "Hopfen_%d_alpha"}
	rm.Hops, err = conv.Hop()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Name": "Hopfen_%d_Sorte", "Amount": "Hopfen_%d_Menge", "Time": "Hopfen_%d_Kochzeit", "Alpha": "Hopfen_%d_alpha"}
	rm.Whirlpool, err = conv.WhirlpoolHop()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Name": "WeitereZutat_Wuerze_%d_Name", "Amount": "WeitereZutat_Wuerze_%d_Menge", "Unit": "WeitereZutat_Wuerze_%d_Einheit", "Time": "WeitereZutat_Wuerze_%d_Kochzeit"}
	rm.Ingredients, err = conv.RecipeTimeUnits()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Name": "WeitereZutat_Gaerung_%d_Name", "Amount": "WeitereZutat_Gaerung_%d_Menge", "Unit": "WeitereZutat_Gaerung_%d_Einheit"}
	rm.Fermentation.Ingredients, err = conv.Ingredient()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Name": "Stopfhopfen_%d_Sorte", "Amount": "Stopfhopfen_%d_Menge"}
	rm.Fermentation.Hops, err = conv.FermentationHop()
	if err != nil {
		return err
	}

	return nil
}

func (con *converter) RecipeUnit() (*recipeUnit, error) {

	k := fmt.Sprintf(con.keys["Name"], con.pos)
	recipeUnit := &recipeUnit{}
	var ok bool
	var err error

	if recipeUnit.Name, ok = con.cmap[k].(string); !ok {
		return nil, errors.New("Key not exists")
	}

	if recipeUnit.Name == "" {
		return nil, fmt.Errorf(fmt.Sprintf("Key Value is empty: %s", k))
	}

	k = fmt.Sprintf(con.keys["Amount"], con.pos)
	if !keyExists(con.cmap, k) {
		return nil, fmt.Errorf(fmt.Sprintf("Amount missing: %s", k))
	}

	switch con.cmap[k].(type) {
	case string:
		if recipeUnit.Amount, err = strconv.ParseFloat(con.cmap[k].(string), 64); err != nil {
			return nil, fmt.Errorf("Parsing Amount error: %s %v", k, err)
		}
	case float64:
		recipeUnit.Amount = con.cmap[k].(float64)
	}

	if _, ok := con.keys["Unit"]; !ok {
		return recipeUnit, nil
	}

	k = fmt.Sprintf(con.keys["Unit"], con.pos)
	if !keyExists(con.cmap, k) {
		return nil, fmt.Errorf(fmt.Sprintf("Unit missing: %s", k))
	}

	if con.cmap[k] == "kg" {
		recipeUnit.Amount = recipeUnit.Amount * 1000
	}

	return recipeUnit, nil
}

func (con *converter) RecipeTimeUnit() (*recipeTimeUnit, error) {
	timeUnit := &recipeTimeUnit{}
	u, err := con.RecipeUnit()
	if err != nil {
		return nil, err
	}
	timeUnit.Name = u.Name
	timeUnit.Amount = u.Amount

	k := fmt.Sprintf(con.keys["Time"], con.pos)
	if !keyExists(con.cmap, k) {
		return nil, fmt.Errorf(fmt.Sprintf("Time missing: %s", k))
	}

	if timeUnit.Time, err = strconv.Atoi(con.cmap[k].(string)); err != nil {
		return nil, err
	}

	return timeUnit, nil
}

func (con *converter) Malts() ([]recipeUnit, error) {
	var Ru []recipeUnit

	for i := 1; keyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		ru, err := con.RecipeUnit()
		if err != nil {
			return nil, err
		}
		Ru = append(Ru, *ru)
	}
	return Ru, nil
}

func (con *converter) RecipeTimeUnits() ([]recipeTimeUnit, error) {
	var Rtu []recipeTimeUnit

	for i := 1; keyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		rtu, err := con.RecipeTimeUnit()
		if err != nil {
			return nil, err
		}
		Rtu = append(Rtu, *rtu)
	}
	return Rtu, nil
}

func (con *converter) Ingredient() ([]recipeTimeUnit, error) {
	var Rtu []recipeTimeUnit

	for i := 1; keyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		rtu := &recipeTimeUnit{}
		ru, err := con.RecipeUnit()
		if err != nil {
			return nil, err
		}
		rtu.Name = ru.Name
		rtu.Amount = ru.Amount
		// this recipe type has no time in this section, use default

		Rtu = append(Rtu, *rtu)
	}
	return Rtu, nil
}

func (con *converter) FermentationHop() ([]Hop, error) {
	var Hops []Hop

	for i := 1; keyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		hop := &Hop{}
		ru, err := con.RecipeUnit()
		if err != nil {
			return nil, err
		}
		hop.Name = ru.Name
		hop.Amount = ru.Amount
		// this recipe type has no time and alpha in this section, use default

		Hops = append(Hops, *hop)
	}

	return Hops, nil
}

func (con *converter) Rests() ([]Rest, error) {
	var Rests []Rest

	for i := 1; keyExists(con.cmap, fmt.Sprintf(con.keys["Time"], i)); i++ {
		con.pos = i
		rest := &Rest{}
		k := fmt.Sprintf(con.keys["Time"], con.pos)
		var ok bool
		var err error

		if _, ok = con.cmap[k].(string); !ok {
			return nil, errors.New("Rests Time Key not exists")
		}

		if con.cmap[k].(string) == "" {
			return nil, fmt.Errorf(fmt.Sprintf("Rest value is empty: %s", k))
		}

		if rest.Time, err = strconv.Atoi(con.cmap[k].(string)); err != nil {
			return nil, fmt.Errorf("Parsing Amount error: %s %v", k, err)
		}

		k = fmt.Sprintf(con.keys["Temperatur"], con.pos)
		if !keyExists(con.cmap, k) {
			return nil, fmt.Errorf(fmt.Sprintf("Amount missing: %s", k))
		}

		if rest.Temperatur, err = strconv.ParseFloat(con.cmap[k].(string), 64); err != nil {
			return nil, fmt.Errorf("Parsing Amount error: %s %v", k, err)
		}

		Rests = append(Rests, *rest)
	}
	return Rests, nil
}

func (con *converter) BasicHop() (*Hop, error) {
	hop := &Hop{}

	u, err := con.RecipeUnit()
	if err != nil {
		return nil, err
	}

	hop.Name = u.Name
	hop.Amount = u.Amount

	k := fmt.Sprintf(con.keys["Alpha"], con.pos)
	if !keyExists(con.cmap, k) {
		return nil, fmt.Errorf(fmt.Sprintf("Alpha missing: %s", k))
	}

	if hop.Alpha, err = strconv.ParseFloat(con.cmap[k].(string), 64); err != nil {
		return nil, fmt.Errorf("Parsing Alpha error: %v", err)
	}
	return hop, nil
}

func (con *converter) FontHop() ([]Hop, error) {
	var Hops []Hop

	for i := 1; keyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		//hop := &Hop{}
		con.pos = i
		hop, err := con.BasicHop()
		if err != nil {
			return nil, err
		}
		Hops = append(Hops, *hop)
	}
	return Hops, nil

}

func (con *converter) Hop() ([]Hop, error) {
	var Hops []Hop

	for i := 1; keyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		hop, err := con.BasicHop()
		if err != nil {
			return nil, err
		}

		k := fmt.Sprintf(con.keys["Time"], i)
		if !keyExists(con.cmap, k) {
			return nil, fmt.Errorf(fmt.Sprintf("Time missing: %s", k))
		}

		if con.cmap[k].(string) == "Whirlpool" {
			continue
		}

		if hop.Time, err = strconv.Atoi(con.cmap[k].(string)); err != nil {
			return nil, err
		}
		Hops = append(Hops, *hop)
	}
	return Hops, nil
}

func (con *converter) WhirlpoolHop() ([]Hop, error) {
	var Hops []Hop
	for i := 1; keyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		hop, err := con.BasicHop()
		if err != nil {
			return nil, err
		}

		k := fmt.Sprintf(con.keys["Time"], i)
		if !keyExists(con.cmap, k) {
			return nil, fmt.Errorf(fmt.Sprintf("Time missing: %s", k))
		}

		if con.cmap[k].(string) == "Whirlpool" {
			Hops = append(Hops, *hop)
		}

	}
	return Hops, nil
}
