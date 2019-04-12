package recipe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tmpdir = filepath.Join(os.TempDir(), "data")

func TestMain(m *testing.M) {
	if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
		os.Mkdir(tmpdir, 0755)
	}
	m.Run()
	os.RemoveAll(tmpdir)
}

func TestRecipeSave(t *testing.T) {
	r, err := LoadFile("testdata/hopsHoney.json", &RecipeM3{})
	assert.Nil(t, err)

	err = r.SavePretty(filepath.Join(tmpdir, "recipeHopsHoney.json"))
	assert.Nil(t, err)
}
func TestRecipeSaveYaml(t *testing.T) {
	r, err := LoadFile("testdata/hopsHoney.json", &RecipeM3{})
	assert.Nil(t, err)

	err = r.SavePrettyYaml(filepath.Join(tmpdir, "recipeHopsHoney.yaml"))
	assert.Nil(t, err)
}

func TestRecipeLoadJson(t *testing.T) {
	_, err := LoadFile(filepath.Join(tmpdir, "recipeHopsHoney.json"), &Recipe{})
	assert.Nil(t, err)
}

func TestRecipeLoadYaml(t *testing.T) {
	_, err := LoadFile(filepath.Join(tmpdir, "recipeHopsHoney.yaml"), &Recipe{})
	assert.Nil(t, err)
}

func TestRecipeSaveJson(t *testing.T) {
	r, err := LoadFile(filepath.Join(tmpdir, "recipeHopsHoney.yaml"), &Recipe{})
	assert.Nil(t, err)
	err = r.SavePretty(filepath.Join(tmpdir, "recipeHopsHoney.json"))
	assert.Nil(t, err)
}
