package recipe

type Malz struct {
	Malz    string  `json:"Malz1"`
	Menge   float32 `json:"Malz1Menge"`
	Einheit string  `json:"Malz1Einheit"`
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

	InfusionHauptguss string `json:"Infusion_Hauptguss"`
	DekoktionEzw      string `json:"Dekoktion_Einmaisch_Zubruehwasser_gesamt"`

	Malz1        string  `json:"Malz1"`
	Malz1Menge   float32 `json:"Malz1_Menge"`
	Malz1Einheit string  `json:"Malz1_Einheit"`

	Malz2        string  `json:"Malz2"`
	Malz2Menge   float32 `json:"Malz2_Menge"`
	Malz2Einheit string  `json:"Malz2_Einheit"`

	Malz3        string  `json:"Malz3"`
	Malz3Menge   float32 `json:"Malz3_Menge"`
	Malz3Einheit string  `json:"Malz3_Einheit"`

	Malz4        string  `json:"Malz4"`
	Malz4Menge   float32 `json:"Malz4_Menge"`
	Malz4Einheit string  `json:"Malz4_Einheit"`

	Malz5        string  `json:"Malz5"`
	Malz5Menge   float32 `json:"Malz5_Menge"`
	Malz5Einheit string  `json:"Malz5_Einheit"`

	Malz6        string  `json:"Malz6"`
	Malz6Menge   float32 `json:"Malz6_Menge"`
	Malz6Einheit string  `json:"Malz6_Einheit"`

	Malz7        string  `json:"Malz7"`
	Malz7Menge   float32 `json:"Malz7_Menge"`
	Malz7Einheit string  `json:"Malz7_Einheit"`
	original     string
}
