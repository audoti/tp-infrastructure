package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/audoti/tp-infrastructure/pkg/fakeiot"
	"github.com/audoti/tp-infrastructure/pkg/putils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	config := putils.GetConfig()
	host := fmt.Sprintf("tcp://%s:%s", config.HOST, config.PORT)
	opts := mqtt.NewClientOptions().AddBroker(host).SetClientID(config.Sensor)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	d := time.Duration(rand.Int31n(config.Delay)) * time.Millisecond
	for x := range time.Tick(d) {
		fmt.Println(x)
		fakeiot.SendWind(c, config)
	}
}
