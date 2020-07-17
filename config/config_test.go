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
