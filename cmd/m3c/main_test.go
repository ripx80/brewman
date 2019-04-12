package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cpath string
var testdata = filepath.Join("..", "..", "pkgs", "recipe", "testdata")
var tmpdir = filepath.Join(os.TempDir(), "data")

func TestMain(m *testing.M) {
	var err error

	cpath, err = os.Getwd()
	if err != nil {
		fmt.Printf("can't get current dir :%s \n", err)
		os.Exit(1)
	}

	cpath = filepath.Join(cpath, "m3c")

	cmd := exec.Command("go", "build", "-o", cpath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("compilation error :%s \n", output)
		os.Exit(1)
	}

	if _, err := os.Stat(tmpdir); os.IsNotExist(err) {
		os.Mkdir(tmpdir, 0755)
	}

	code := m.Run()
	os.Remove(cpath)
	os.RemoveAll(tmpdir)
	os.Exit(code)
}

func TestConvert(t *testing.T) {
	m3c := exec.Command(cpath, filepath.Join(testdata, "apiTest.json"), filepath.Join(tmpdir, "test.json"))
	err := m3c.Run()
	assert.Nil(t, err)
}

func TestConvertYaml(t *testing.T) {
	m3c := exec.Command(cpath, "--o", "yaml", filepath.Join(testdata, "apiTest.json"), filepath.Join(tmpdir, "test.yaml"))
	err := m3c.Run()
	assert.Nil(t, err)
}
