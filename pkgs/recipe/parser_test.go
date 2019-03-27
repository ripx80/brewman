package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiRecipe(t *testing.T) {
	_, err := LoadFile("testdata/test.json")
	assert.Nil(t, err)
}
