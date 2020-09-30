package main

import (
"github.com/tkanos/gonfig"
"os"
"path/filepath"
)

type Configuration struct {
	Delay int
	Type string
	Captor string
	AITA string
	TOPIC string
	PORT string
	HOST string
}

func GetConfig () Configuration {

	configuration := Configuration{}

	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	err = gonfig.GetConf(filepath.Join(dir, "aeroport", "config", "wind_config.json"), &configuration)

	if err != nil {
		panic(err)
	}

	return configuration
	
}