package putils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	newConfig := Configuration{
		Delay:  delay,
		Type:   sensorType,
		Sensor: sensorID,
		AITA:   aita,
		TOPIC:  topic,
		PORT:   port,
		HOST:   host,
	}
	fmt.Println(newConfig)
	return newConfig
}

// TimeToDate will convert a golang `time.Now()` to a YYYY-MM-DD date
func TimeToDate(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

// Sleep <ms> milliseconds
func Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// Aita is a list of airport codes
var Aita = []string{"AMS", "ATL", "BKK", "CAN", "CDG", "DEL", "DEN", "DFW", "DXB", "FRA", "HKG", "HND", "ICN", "JFK", "LAX", "LHR", "ORD", "PEK", "PVG", "SIN"}

// Airport is an airport's data
type Airport struct {
	Aita string
	Name string
}

// AitaFull is a list of airport codes with its full name
var AitaFull = []Airport{
	Airport{"AMS", "Amsterdam Airport Schiphol"},
	Airport{"ATL", "Hartsfield–Jackson Atlanta International Airport"},
	Airport{"BKK", "Suvarnabhumi Airport"},
	Airport{"CAN", "Guangzhou Baiyun International Airport"},
	Airport{"CDG", "Paris Charles de Gaulle Airport"},
	Airport{"DEL", "Indira Gandhi International Airport"},
	Airport{"DEN", "Denver International Airport"},
	Airport{"DFW", "Dallas/Fort Worth International Airport"},
	Airport{"DXB", "Dubai International Airport"},
	Airport{"FRA", "Frankfurt am Main Airport"},
	Airport{"HKG", "Hong Kong International Airport"},
	Airport{"HND", "Tokyo International Airport"},
	Airport{"ICN", "Incheon International Airport"},
	Airport{"JFK", "John F. Kennedy International Airport"},
	Airport{"LAX", "Los Angeles International Airport"},
	Airport{"LHR", "Heathrow Airport"},
	Airport{"ORD", "O'Hare International Airport"},
	Airport{"PEK", "Beijing Capital International Airport"},
	Airport{"PVG", "Shanghai Pudong International Airport"},
	Airport{"SIN", "Singapore Changi Airport"},
}
