package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//var tmpdir = filepath.Join(os.TempDir(), "data")

// func TestMain(m *testing.M) {
// 	if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
// 		os.Mkdir(tmpdir, 0755)
// 	}
// 	m.Run()
// 	fmt.Println("remove tmpdir")
// 	os.RemoveAll(tmpdir)
// }
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
	_, err := LoadFile("testdata/hopsHoney.json", &RecipeM3{})
	assert.Nil(t, err)
}
