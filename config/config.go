package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// type PodsConfig struct {
// 	Hotwater PodConfig `yaml:"hotwater" validate:"min=0,max=40"`
// 	Masher   PodConfig `yaml:"masher" validate:"min=0,max=40"`
// 	Cooker   PodConfig `yaml:"cooker" validate:"min=0,max=40"`
// }

type TemperaturConfig struct {
	Device  string `yaml:"device"` // only ds18b20 supp.
	Bus     int    `yaml:"bus"`
	Address int    `yaml:"address"`
}

type PodConfig struct {
	Control    int `yaml:"control" validate:"min=0,max=40"`
	Agiator    int `yaml:"agiator" validate:"min=0,max=40"`
	Temperatur TemperaturConfig
}
type Config struct {
	Global   GlobalConfig `yaml:"global"`
	Hotwater PodConfig    `yaml:"hotwater"`
	Masher   PodConfig    `yaml:"masher"`
	Cooker   PodConfig    `yaml:"cooker"`
	//Sensor  SensorConfig  `yaml:"sensors"`
	//Control ControlConfig `yaml:"control"`
	Recipe RecipeConfig `yaml:"recipe"`
	// original is the input from which the config was parsed.
	original string
}

type GlobalConfig struct {
	TemperaturUnit string `yaml:"temperatur-unit" validate:"regexp=[Cc]elsius|[Kk]elvin"`
}

// type SensorConfig struct {
// 	Hotwater int `yaml:"hotwater" validate:"min=0,max=40"`
// 	Masher   int `yaml:"masher" validate:"min=0,max=40"`
// 	Cooker   int `yaml:"cooker" validate:"min=0,max=40"`
// 	FlowIn   int `yaml:"flowin" validate:"min=0,max=40"`
// }

// type ControlConfig struct {
// 	HeaterWater  int `yaml:"heater-water" validate:"min=0,max=40"`
// 	HeaterMash   int `yaml:"heater-mash" validate:"min=0,max=40"`
// 	HeaterCooker int `yaml:"heater-cooker" validate:"min=0,max=40"`
// 	//Gpio []map[string]string `yaml:"gpio"`
// }

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
		TemperaturUnit: "Celsius",
	}

	// DefaultPodConfig = PodConfig{
	// 	Control: 10,
	// }

	// DefaultControlConfig = ControlConfig{
	// 	HeaterWater:  29,
	// 	HeaterMash:   31,
	// 	HeaterCooker: 32,
	// }

	DefaultRecipeConfig = RecipeConfig{
		File: "recipe.yaml",
	}
)

func Load(s string) (*Config, error) {
	cfg := &Config{}
	//init default config
	*cfg = DefaultConfig

	err := yaml.UnmarshalStrict([]byte(s), cfg)
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
