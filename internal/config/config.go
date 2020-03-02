package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/rsmaxwell/page/internal/myfile"
)

// Config structure
type Config struct {
	DocumentRoot string `json:"documentRoot"`
}

// New creates a config object
func New() *Config {
	c := new(Config)

	c.DocumentRoot = "/var/www/"

	configfile := "/etc/page/page.json"

	if runtime.GOOS == "windows" {
		c.DocumentRoot = "C:/temp"
		configfile = "C:/temp/page.json"
	}

	value, ok := os.LookupEnv("PAGE_CONFIGFILE")
	if ok {
		configfile = value
	}

	if myfile.Exists(configfile) {
		jsonFile, err := os.Open(configfile)
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
