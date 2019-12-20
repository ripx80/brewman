package recipe

import (
	"fmt"
	"io/ioutil"
)

/*
Parser interface to Load a file and convert it
*/
type Parser interface {
	Load(s string) (*Recipe, error)
	String() string
	PrettyPrint() string
}

func keyExists(m map[string]interface{}, k string) bool {
	if _, ok := m[k]; ok {
		return true
	}
	return false
}

/*
LoadFile file from disk and convert
*/
func LoadFile(filename string, p Parser) (*Recipe, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	recipe, err := p.Load(string(content))

	if err != nil {
		return nil, fmt.Errorf("parsing recipe file %s: %v", filename, err)
	}
	return recipe, nil
}
