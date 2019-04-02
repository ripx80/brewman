package recipe

import (
	"fmt"
	"io/ioutil"
)

type Parser interface {
	Load(s string) (*Recipe, error)
	String() string
	PrettyPrint() string
}

func KeyExists(m map[string]interface{}, k string) bool {
	if _, ok := m[k]; ok {
		return true
	}
	return false
}

func LoadFile(filename string, p Parser) (*Recipe, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	//m3 := RecipeM3{}
	recipe, err := p.Load(string(content))

	if err != nil {
		return nil, fmt.Errorf("parsing recipe file %s: %v", filename, err)
	}
	return recipe, nil
}
