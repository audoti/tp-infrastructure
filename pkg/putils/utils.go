package putils

import (
	"os"
	"path/filepath"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Delay  int32
	Type   string
	Sensor string
	AITA   string
	TOPIC  string
	PORT   string
	HOST   string
}

func GetConfig() Configuration {
	configuration := Configuration{}
	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	err = gonfig.GetConf(filepath.Join(dir, "configs", "fakeiot_wind_config.json"), &configuration)

	if err != nil {
		panic(err)
	}

	return configuration
}
