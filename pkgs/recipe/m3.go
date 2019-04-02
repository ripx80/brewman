package recipe

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Converter struct {
	keys map[string]string
	cmap map[string]interface{}
	pos  int
}

type RecipeM3Global struct {
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
}

type RecipeM3 struct {
	RecipeM3Global

	MainCast      float64 `json:"Infusion_Hauptguss,string"`
	InTemperatur  float64 `json:"Infusion_Einmaischtemperatur"`
	OutTemperatur float64 `json:"Abmaischtemperatur,string"`
	Grouting      float64 `json:"Nachguss,string"`
	Time          int     `json:"Kochzeit_Wuerze,string"`
	Yeast         string  `json:"Hefe"`
	Temperatur    float64 `json:"Gaertemperatur,string"` // json field values has 24-25 as string, not supported
	EndDegree     float64 `json:"Endvergaerungsgrad,string"`
	Carbonation   float64 `json:"Karbonisierung,string"`
	Annotation    string  `json:"Anmerkung_Autor"`

	Malts        []Malt
	Rests        []Rest
	FontHops     []Hop
	Hops         []Hop
	Whirpool     []Hop
	Ingredients  []Ingredient
	Fermentation RecipeFermentation
}

func (rm *RecipeM3) Load(s string) (*Recipe, error) {
	err := json.Unmarshal([]byte(s), rm)
	if err != nil {
		return nil, err
	}
	recipe := &Recipe{}
	//recipe.original = s
	recipe.Global = RecipeGlobal{
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
	}

	recipe.Water = RecipeWater{
		MainCast: rm.MainCast,
		Grouting: rm.Grouting,
	}

	recipe.Mash = RecipeMash{
		InTemperatur:  rm.InTemperatur,
		OutTemperatur: rm.OutTemperatur,
		Malts:         rm.Malts,
	}

	recipe.Cook = RecipeCook{
		Time:        rm.Time,
		Ingredients: rm.Ingredients,
		FontHops:    rm.FontHops,
		Hops:        rm.Hops,
		Whirpool:    rm.Whirpool,
	}

	recipe.Fermentation = RecipeFermentation{
		Yeast:       rm.Yeast,
		Temperatur:  rm.Temperatur,
		EndDegree:   rm.EndDegree,
		Carbonation: rm.Carbonation,
		Hops:        rm.Fermentation.Hops,
		Ingredients: rm.Fermentation.Ingredients,
	}
	return recipe, nil
}

func (rm RecipeM3) String() string {
	b, err := json.Marshal(rm)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (rm RecipeM3) PrettyPrint() string {
	b, err := json.MarshalIndent(rm, "", "   ")
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (rm *RecipeM3) UnmarshalJSON(data []byte) error {
	// this is no good implementation. I hope my skills will be better in the future to do this -.-
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

	conv := &Converter{}
	conv.cmap = result
	//Todo: convert Einheit kg to g
	conv.keys = map[string]string{"Name": "Malz%d", "Amount": "Malz%d_Menge"}
	rm.Malts, err = conv.RecipeUnits()
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
	rm.Whirpool, err = conv.WhirpoolHop()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Name": "WeitereZutat_Wuerze_%d_Name", "Amount": "WeitereZutat_Wuerze_%d_Menge", "Time": "WeitereZutat_Wuerze_%d_Kochzeit"}
	rm.Ingredients, err = conv.RecipeTimeUnits()
	if err != nil {
		return err
	}

	conv.keys = map[string]string{"Name": "WeitereZutat_Gaerung_%d_Name", "Amount": "WeitereZutat_Gaerung_%d_Menge"}
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

func (con *Converter) RecipeUnit() (*RecipeUnit, error) {

	k := fmt.Sprintf(con.keys["Name"], con.pos)
	recipeUnit := &RecipeUnit{}
	var ok bool
	var err error

	if recipeUnit.Name, ok = con.cmap[k].(string); !ok {
		return nil, errors.New("Key not exists")
	}

	if recipeUnit.Name == "" {
		return nil, errors.New(fmt.Sprintf("Key Value is empty: %s", k))
	}

	k = fmt.Sprintf(con.keys["Amount"], con.pos)
	if !KeyExists(con.cmap, k) {
		return nil, errors.New(fmt.Sprintf("Amount missing: %s", k))
	}

	switch con.cmap[k].(type) {
	case string:
		if recipeUnit.Amount, err = strconv.ParseFloat(con.cmap[k].(string), 64); err != nil {
			return nil, fmt.Errorf("Parsing Amount error: %s %v", k, err)
		}
	case float64:
		recipeUnit.Amount = con.cmap[k].(float64)
	}

	return recipeUnit, nil
}

func (con *Converter) RecipeTimeUnit() (*RecipeTimeUnit, error) {
	recipeTimeUnit := &RecipeTimeUnit{}
	u, err := con.RecipeUnit()
	if err != nil {
		return nil, err
	}
	recipeTimeUnit.Name = u.Name
	recipeTimeUnit.Amount = u.Amount

	k := fmt.Sprintf(con.keys["Time"], con.pos)
	if !KeyExists(con.cmap, k) {
		return nil, errors.New(fmt.Sprintf("Time missing: %s", k))
	}

	if recipeTimeUnit.Time, err = strconv.Atoi(con.cmap[k].(string)); err != nil {
		return nil, err
	}

	return recipeTimeUnit, nil
}

func (con *Converter) RecipeUnits() ([]RecipeUnit, error) {
	var Ru []RecipeUnit

	for i := 1; KeyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		ru, err := con.RecipeUnit()
		if err != nil {
			return nil, err
		}
		Ru = append(Ru, *ru)
	}
	return Ru, nil
}

func (con *Converter) RecipeTimeUnits() ([]RecipeTimeUnit, error) {
	var Rtu []RecipeTimeUnit

	for i := 1; KeyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		rtu, err := con.RecipeTimeUnit()
		if err != nil {
			return nil, err
		}
		Rtu = append(Rtu, *rtu)
	}
	return Rtu, nil
}

func (con *Converter) Ingredient() ([]RecipeTimeUnit, error) {
	var Rtu []RecipeTimeUnit

	for i := 1; KeyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		rtu := &RecipeTimeUnit{}
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

func (con *Converter) FermentationHop() ([]Hop, error) {
	var Hops []Hop

	for i := 1; KeyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
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

func (con *Converter) Rests() ([]Rest, error) {
	var Rests []Rest

	for i := 1; KeyExists(con.cmap, fmt.Sprintf(con.keys["Time"], i)); i++ {
		con.pos = i
		rest := &Rest{}
		k := fmt.Sprintf(con.keys["Time"], con.pos)
		var ok bool
		var err error

		if _, ok = con.cmap[k].(string); !ok {
			return nil, errors.New("Rests Time Key not exists")
		}

		if con.cmap[k].(string) == "" {
			return nil, errors.New(fmt.Sprintf("Rest value is empty: %s", k))
		}

		if rest.Time, err = strconv.Atoi(con.cmap[k].(string)); err != nil {
			return nil, fmt.Errorf("Parsing Amount error: %s %v", k, err)
		}

		k = fmt.Sprintf(con.keys["Temperatur"], con.pos)
		if !KeyExists(con.cmap, k) {
			return nil, errors.New(fmt.Sprintf("Amount missing: %s", k))
		}

		if rest.Temperatur, err = strconv.ParseFloat(con.cmap[k].(string), 64); err != nil {
			return nil, fmt.Errorf("Parsing Amount error: %s %v", k, err)
		}

		Rests = append(Rests, *rest)
	}
	return Rests, nil
}

func (con *Converter) BasicHop() (*Hop, error) {
	hop := &Hop{}

	u, err := con.RecipeUnit()
	if err != nil {
		return nil, err
	}

	hop.Name = u.Name
	hop.Amount = u.Amount

	k := fmt.Sprintf(con.keys["Alpha"], con.pos)
	if !KeyExists(con.cmap, k) {
		return nil, errors.New(fmt.Sprintf("Alpha missing: %s", k))
	}

	if hop.Alpha, err = strconv.ParseFloat(con.cmap[k].(string), 64); err != nil {
		return nil, fmt.Errorf("Parsing Alpha error: %v", err)
	}
	return hop, nil
}

func (con *Converter) FontHop() ([]Hop, error) {
	var Hops []Hop

	for i := 1; KeyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
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

func (con *Converter) Hop() ([]Hop, error) {
	var Hops []Hop

	for i := 1; KeyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		hop, err := con.BasicHop()
		if err != nil {
			return nil, err
		}

		k := fmt.Sprintf(con.keys["Time"], i)
		if !KeyExists(con.cmap, k) {
			return nil, errors.New(fmt.Sprintf("Time missing: %s", k))
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

func (con *Converter) WhirpoolHop() ([]Hop, error) {
	var Hops []Hop
	for i := 1; KeyExists(con.cmap, fmt.Sprintf(con.keys["Name"], i)); i++ {
		con.pos = i
		hop, err := con.BasicHop()
		if err != nil {
			return nil, err
		}

		k := fmt.Sprintf(con.keys["Time"], i)
		if !KeyExists(con.cmap, k) {
			return nil, errors.New(fmt.Sprintf("Time missing: %s", k))
		}

		if con.cmap[k].(string) == "Whirlpool" {
			Hops = append(Hops, *hop)
		}

	}
	return Hops, nil
}
