package recipe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeM3Api(t *testing.T) {
	_, err := LoadFile("testdata/apiTest.json", &RecipeM3{})
	assert.Nil(t, err)
}

func TestRecipeM3FermentationHop(t *testing.T) {
	_, err := LoadFile("testdata/fermentationHop.json", &RecipeM3{})
	assert.Nil(t, err)
}

func TestRecipeM3Whirlpool(t *testing.T) {
	_, err := LoadFile("testdata/whirlpool.json", &RecipeM3{})
	assert.Nil(t, err)
}

func TestRecipeM3HopsHoney(t *testing.T) {
	r, err := LoadFile("testdata/hopsHoney.json", &RecipeM3{})
	fmt.Println(r.PrettyPrint())
	assert.Nil(t, err)
}
