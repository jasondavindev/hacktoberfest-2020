package main

import (
	"bytes"
	"io/ioutil"

	"go.uber.org/config"
)

//Config is a struct that will be used to parse the enviroment variables located in config.yml
type Config struct {
	Goenv     string
	Directory string
	Exclude   string
}

//CfgFactory generates config based on config.yml file
func CfgFactory() Config {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	yamlFileReader := bytes.NewReader(yamlFile)

	provider, err := config.NewYAML(config.Source(yamlFileReader))
	if err != nil {
		panic(err)
	}

	var c Config
	if err := provider.Get("development").Populate(&c); err != nil {
		log.Fatal(err)
	}

	return c
}
