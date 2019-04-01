package recipe

type Converter struct {
	keys []string
	cmap map[string]interface{}
	pos  int
}
type RecipeM3 struct {
	Name     string `json:"Name" validate:"nonzero"`
	Datum    string `json:"Datum" validate:"nonzero"`
	Sorte    string `json:"Sorte" validate:"nonzero"`
	Autor    string `json:"Autor" validate:"nonzero"`
	Klonbier string `json:"Klonbier" validate:"nonzero"`

	KlonbierOriginal string  `json:"Klonbier_Original"`
	Ausschlagswuerze int     `json:"Ausschlagswuerze"`
	Sudhausausbeute  int     `json:"Sudhausausbeute"`
	Stammwuerze      float32 `json:"Stammwuerze"`
	Bittere          int     `json:"Bittere"`
	Farbe            string  `json:"Farbe"`
	Alkohol          float32 `json:"Alkohol"`
	Kurzbeschreibung string  `json:"Kurzbeschreibung"`

	Hauptguss           string `json:"Infusion_Hauptguss"`
	Einmaischtemperatur int    `json:"Infusion_Einmaischtemperatur"`
	Abmaischtemperatur  string `json:"Abmaischtemperatur"`
	Nachguss            string `json:"Nachguss"`
	KochzeitWuerze      string `json:"Kochzeit_Wuerze"`
	Hefe                string `json:"Hefe"`
	Gaertemperatur      string `json:"Gaertemperatur"`
	Endvergaerungsgrad  string `json:"Endvergaerungsgrad"`
	Karbonisierung      string `json:"Karbonisierung"`
	AnmerkungAutor      string `json:"Anmerkung_Autor"`

	Malts        []Malt
	Rests        []Rest
	FontHops     []Hop
	Hops         []Hop
	Whirpool     []Hop
	Ingredients  []Ingredient
	Fermentation RecipeFermentation

	// Malts will be directly convert to internal struct
	// Rasts will be directly convert to internal struct

	original string
}
