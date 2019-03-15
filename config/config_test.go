package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {

	_, err := LoadFile("testdata/conf.yaml")
	assert.Nil(t, err)

	// expectedConf.original = c.original
	// assert.Equal(t, expectedConf, c)
}

func TestBadConfigs(t *testing.T) {
	_, err := LoadFile("testdata/conf_bad_sensors.yaml")
	assert.Error(t, err)

}
