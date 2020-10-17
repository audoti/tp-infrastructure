package fakeiot

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/audoti/tp-infrastructure/pkg/putils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func sendWind(c mqtt.Client, config putils.Configuration) {
	mesure := generateMesure(config)
	c.Publish(config.TOPIC, 1, false, mesure)

	if token := c.Publish(config.TOPIC, 1, false, mesure); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	fmt.Println(mesure)
}

func generateMesure(config putils.Configuration) string {
	value := strconv.Itoa(rand.Intn(100))
	date := time.Now().Format(time.RFC3339)
	return fmt.Sprintf("%s,%s,%s,%s,%s", config.Sensor, config.AITA, config.Type, value, date)
}

func StartFakeIot() {
	delay := flag.Int("delay", 100, "the delay of the sensor")
	sensorType := flag.String("sensorType", "Wind speed", "the type of the sensor")
	sensorID := flag.String("sensorID", "0000", "the sensor id")
	aita := flag.String("aita", putils.Aita[0], "the code of the airport the sensor is associated to")
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
	for {
		putils.Sleep(int(d.Milliseconds()))
		sendWind(c, config)
	}
}
