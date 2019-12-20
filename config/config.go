package config

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// add a general type for all. add validator func for supported devices
type TemperaturConfig struct {
	Device  string `yaml:"device" validate:"nonzero"`  // only ds18b20 supp.
	Address string `yaml:"address" validate:"nonzero"` // path to sys file. its bad i know
}

type ControlConfig struct {
	Device  string `yaml:"device" validate:"nonzero"`
	Address string `yaml:"address" validate:"min=0,max=40,nonzero"`
}

type AgiatorConfig struct {
	Device  string `yaml:"device"`
	Address string `yaml:"address" validate:"min=0,max=40`
}

type PodConfig struct {
	Control    ControlConfig
	Agiator    AgiatorConfig
	Temperatur TemperaturConfig
}
type Config struct {
	Global   GlobalConfig `yaml:"global"`
	Hotwater PodConfig    `yaml:"hotwater"`
	Masher   PodConfig    `yaml:"masher"`
	Cooker   PodConfig    `yaml:"cooker"`
	Recipe   RecipeConfig `yaml:"recipe"`
	// original is the input from which the config was parsed.
	original string
}

type GlobalConfig struct {
	TemperaturUnit     string  `yaml:"temperatur-unit" validate:"regexp=[Cc]elsius|[Kk]elvin"`
	HotwaterTemperatur float64 `yaml:"hotwater-temperatur" validate:"min=70,max=90"`
	CookingTemperatur  float64 `yaml:"cooking-temperatur" validate:"min=90,max=110"`
}

type RecipeConfig struct {
	File string `yaml:"file"`
}

var (
	DefaultConfig = Config{
		Global: DefaultGlobalConfig,
		//Hotwater: DefaultPodConfig,
		//Masher:   DefaultPodConfig,
		//Cooker:   DefaultPodConfig,
		Recipe: DefaultRecipeConfig,
	}

	DefaultGlobalConfig = GlobalConfig{
		TemperaturUnit:     "Celsius",
		HotwaterTemperatur: 76.0,
		CookingTemperatur:  97.5,
	}

	// DefaultPodConfig = PodConfig{
	// 	Control: 10,
	// }
	DefaultRecipeConfig = RecipeConfig{
		File: "recipe.yaml",
	}
)

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

// func (cc ControlConfig) String() string {
// 	b, err := yaml.Marshal(cc)
// 	if err != nil {
// 		return fmt.Sprintf("<error creating control config string: %s>", err)
// 	}
// 	return string(b)
// }

func (rc RecipeConfig) String() string {
	b, err := yaml.Marshal(rc)
	if err != nil {
		return fmt.Sprintf("<error creating recipe config string: %s>", err)
	}
	return string(b)
}

// Impmement this interface allows you to parse the config file!

// func (c Config) UnmarshalYAML(unmarshal func(interface{}) error) error{
// 	*c = DefaultConfig
// 	type plain Config
// 	if err:=unmarshal((*plain)(c)); err != nil{
// 		return err
// 	}

// }
