package putils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Delay  int
	Type   string
	Sensor string
	AITA   string
	TOPIC  string
	PORT   int
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

func SetConfig(delay int, sensorType string, sensorID string, aita string, topic string, host string, port int) Configuration {

	newConfif := Configuration{
		Delay:  delay,
		Type:   sensorType,
		Sensor: sensorID,
		AITA:   aita,
		TOPIC:  topic,
		PORT:   port,
		HOST:   host,
	}
	fmt.Println(newConfif)
	return newConfif
}
