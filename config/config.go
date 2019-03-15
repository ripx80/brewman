package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Sensor  SensorConfig  `yaml:"sensors"`
	Control ControlConfig `yaml:"control"`
	// original is the input from which the config was parsed.
	original string
}

type SensorConfig struct {
	TemperaturUnit string `yaml:"temperatur-unit" validate:"regexp=[Cc]elsius|[Kk]elvin"`
	Hotwater       int    `yaml:"hotwater" validate:"min=0,max=40"`
	Masher         int    `yaml:"masher" validate:"min=0,max=40"`
	Cooker         int    `yaml:"cooker" validate:"min=0,max=40"`
	FlowIn         int    `yaml:"flowin" validate:"min=0,max=40"`
}

type ControlConfig struct {
	Gpio []map[string]string `yaml:"gpio"`
}

func Load(s string) (*Config, error) {
	cfg := &Config{}
	//init default config
	//*cfg = DefaultConfig

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
