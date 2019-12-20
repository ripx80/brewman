package config

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

/*
PeriphConfig struct stores the device type and the Address in a string.
*/
type periphConfig struct {
	Device  string `yaml:"device" validate:"nonzero"`
	Address string `yaml:"address" validate:"nonzero"` // path to sys file. its bad i know
}

type periphConfigZero struct {
	Device  string `yaml:"device"`
	Address string `yaml:"address"`
}

/*
PodConfig holds Pod Informations
*/
type PodConfig struct {
	Control    periphConfig
	Agiator    periphConfigZero
	Temperatur periphConfig
}

/*
Config holds all stages of config data
*/
type Config struct {
	Global   GlobalConfig `yaml:"global"`
	Hotwater PodConfig    `yaml:"hotwater"`
	Masher   PodConfig    `yaml:"masher"`
	Cooker   PodConfig    `yaml:"cooker"`
	Recipe   RecipeConfig `yaml:"recipe"`
	// original is the input from which the config was parsed.
	original string
}

/*
GlobalConfig general and global configs
*/
type GlobalConfig struct {
	TemperaturUnit     string  `yaml:"temperatur-unit" validate:"regexp=[Cc]elsius|[Kk]elvin"`
	HotwaterTemperatur float64 `yaml:"hotwater-temperatur" validate:"min=70,max=90"`
	CookingTemperatur  float64 `yaml:"cooking-temperatur" validate:"min=90,max=110"`
}

/*
RecipeConfig configurations for recipes
*/
type RecipeConfig struct {
	File string `yaml:"file"`
}

var (
	/*DefaultConfig holds a basic default configuration*/
	DefaultConfig = Config{
		Global: defaultGlobalConfig,
		//Hotwater: DefaultPodConfig,
		//Masher:   DefaultPodConfig,
		//Cooker:   DefaultPodConfig,
		Recipe: defaultRecipeConfig,
	}

	defaultGlobalConfig = GlobalConfig{
		TemperaturUnit:     "Celsius",
		HotwaterTemperatur: 76.0,
		CookingTemperatur:  97.5,
	}

	// DefaultPodConfig = PodConfig{
	// 	Control: 10,
	// }
	defaultRecipeConfig = RecipeConfig{
		File: "recipe.yaml",
	}
)

/*
Load wrapper func to load a config in yaml format
*/
func Load(s string) (*Config, error) {
	cfg := &Config{}
	//init default config
	*cfg = DefaultConfig

	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	cfg.original = s
	return cfg, nil
}

/*
LoadFile reads the config file and parse the yaml
*/
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing YAML file %s: %v", filename, err)
	}
	return cfg, nil
}

func (c Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

/*
Save config file to disk
*/
func (c Config) Save(fn string) error {
	return ioutil.WriteFile(fn, []byte(c.String()), 0644)
}

func (pc PodConfig) String() string {
	b, err := yaml.Marshal(pc)
	if err != nil {
		return fmt.Sprintf("<error creating sensor config string: %s>", err)
	}
	return string(b)
}

func (rc RecipeConfig) String() string {
	b, err := yaml.Marshal(rc)
	if err != nil {
		return fmt.Sprintf("<error creating recipe config string: %s>", err)
	}
	return string(b)
}
