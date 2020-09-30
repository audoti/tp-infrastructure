package main

import "github.com/eclipse/paho.mqtt.golang"
import "fmt"
import "time"
import "strconv"
import "math/rand"

func main() {

	config := GetConfig()

	host := fmt.Sprintf("tcp://%s:%s", config.HOST, config.PORT)

	opts := mqtt.NewClientOptions().AddBroker(host).SetClientID(config.Captor)

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	d := time.Duration(rand.Int31n(config.Delay)) *  time.Millisecond

	for x := range time.Tick(d) {
		fmt.Println(x)
		SendWind(c, config)
	}

}

func SendWind(c mqtt.Client, config Configuration) {
			
	mesure := GenerateMesure(config)

	c.Publish(config.TOPIC, 1, false, mesure)

	if token := c.Publish(config.TOPIC, 1, false, mesure); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	
	fmt.Println(mesure)
}

func GenerateMesure (config Configuration) string {
	value := strconv.Itoa(rand.Intn(100))
	date := time.Now()
	return fmt.Sprintf("%s:%s:%s:%s:%s",config.Captor, config.AITA, config.Type, value, date)
}