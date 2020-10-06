package fakeiot

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/audoti/tp-infrastructure/pkg/putils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SendWind(c mqtt.Client, config putils.Configuration) {
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
