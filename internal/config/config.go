package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config is the Config structure
type Config struct {
	Prefix string `json:"prefix"`
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// New creates a config object
func New() *Config {
	c := new(Config)

	filename := "/etc/page/page.json"
	if fileExists(filename) {
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
