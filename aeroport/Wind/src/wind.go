package main


import "github.com/eclipse/paho.mqtt.golang"
import "fmt"
import "time"

func main() {

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("aeroport_wind")

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	d := 200*time.Millisecond

	for x := range time.Tick(d) {
		fmt.Println(x)
		sendWind(c)
	}

}

func sendWind(c mqtt.Client) {
			
	c.Publish("aeroport/wind", 1, false, "20KM/h")

	if token := c.Publish("aeroport/wind", 1, false, "20KM/h"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	
	fmt.Println("Data sent")
}
