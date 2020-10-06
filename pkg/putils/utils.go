package putils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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

// TimeToDate will convert a golang `time.Now()` to a YYYY-MM-DD date
func TimeToDate(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
