package recipe

// all units in g, l

type Recipe struct {
	Global       RecipeGlobal       `json:"Global"`
	Water        RecipeWater        `json:"Water"`
	Mash         RecipeMash         `json:"Mash"`
	Cook         RecipeCook         `json:"Cook"`
	Fermentation RecipeFermentation `json:"Fermentation"`
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

type Malt struct {
	Name   string
	Amount float64
}

type Rest struct {
	Time       int
	Temperatur float64
}

type Hop struct {
	Name   string
	Amount int
	Alpha  float64
}

type Ingredient struct {
	Name   string
	Amount float64 //in g
}
