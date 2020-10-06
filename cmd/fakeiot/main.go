package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/audoti/tp-infrastructure/pkg/fakeiot"
	"github.com/audoti/tp-infrastructure/pkg/putils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	delay := flag.Int("delay", 100, "the delay of the sensor")
	sensorType := flag.String("sensorType", "Wind speed", "the type of the sensor")
	sensorID := flag.String("sensorID", "0000", "the sensor id")
	aita := flag.String("aita", "AITA", "the code of the airport the sensor is associated to")
	topic := flag.String("topic", "aeroport/wind", "the topic of the sensor")
	sensorHost := flag.String("host", "localhost", "the host")
	port := flag.Int("port", 18830, "the port")
	flag.Parse()

	config := putils.SetConfig(*delay, *sensorType, *sensorID, *aita, *topic, *sensorHost, *port)
	//config := putils.GetConfig()
	host := fmt.Sprintf("tcp://%s:%d", config.HOST, config.PORT)
	opts := mqtt.NewClientOptions().AddBroker(host).SetClientID(config.Sensor)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	d := time.Duration(rand.Int31n(int32(config.Delay))) * time.Millisecond
	for x := range time.Tick(d) {
		fmt.Println(x)
		fakeiot.SendWind(c, config)
	}
}
