package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rsmaxwell/page/internal/file"
)

// Config structure
type Config struct {
	Prefix string `json:"prefix"`
	Debug  Debug  `json:"debug"`
}

// Debug structure
type Debug struct {
	Directory            string         `json:"directory"`
	Filename             string         `json:"filename"`
	Level                int            `json:"level"`
	DefaultPackageLevel  int            `json:"defaultPackageLevel"`
	DefaultFunctionLevel int            `json:"defaultFunctionLevel"`
	DumpDir              string         `json:"dumpDir"`
	FunctionLevels       map[string]int `json:"functionLevels"`
	PackageLevels        map[string]int `json:"packageLevels"`
}

// New creates a config object
func New() *Config {
	c := new(Config)

	c.Prefix = "/var/www/"
	c.Debug.Filename = "/tmp/page.log"
	c.Debug.Level = 30
	c.Debug.DefaultPackageLevel = 30
	c.Debug.DefaultFunctionLevel = 30
	c.Debug.DumpDir = "/tmp"
	c.Debug.FunctionLevels = make(map[string]int)
	c.Debug.PackageLevels = make(map[string]int)

	filename := "/etc/page/page.json"
	if file.Exists(filename) {
		jsonFile, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(byteValue, c)
	}

	return c
}
