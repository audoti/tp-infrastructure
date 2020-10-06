package mqttfile

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/audoti/tp-infrastructure/pkg/putils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// StartMqttFile will start the mqtt listener and write to a csv file
func StartMqttFile() {
	client := connect("tcp://localhost:18830", "listener-file")
	subscribeMqtt(client)
	// Kill the process on SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func createClientOptions(brokerURI string, clientID string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	// AddBroker adds a broker URI to the list of brokers to be used.
	// The format should be "scheme://host:port"
	opts.AddBroker(brokerURI)
	// opts.SetUsername(user)
	// opts.SetPassword(password)
	opts.SetClientID(clientID)
	return opts
}

func connect(brokerURI string, clientID string) mqtt.Client {
	fmt.Println("Trying to connect (" + brokerURI + ", " + clientID + ")...")
	opts := createClientOptions(brokerURI, clientID)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func subscribeMqtt(client mqtt.Client) {
	// client.Subscribe("aeroport/#", 1, onMessageReceived)
	topic := "aeroport/#"
	if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("Connected to MQTT broker.")
		fmt.Println("Subscribed to the topic \"" + topic + "\". Reacting to incoming messages.")
	}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	// stock le contenu du message
	body := string(message.Payload())
	// stock le topic du message
	topic := message.Topic()
	// stock le capteur
	typeSensor := strings.Split(topic, "/")[1]
	// split le body puis récupère les informations envoyés
	splitted := strings.Split(body, ",")
	idSensor := splitted[0]
	airportCode := splitted[1]
	valueSensor := splitted[3]
	dateTime, _ := time.Parse(time.RFC3339, splitted[4])
	// parse la dateTime en date
	date := putils.TimeToDate(dateTime)
	// chemin du fichier csv
	filePath := filepath.Join("logs", airportCode+"-"+date+"-"+typeSensor+".csv")
	fmt.Println(filePath)
	writeFile(filePath, fmt.Sprintf("%s,%s,%s,%s,%s", idSensor, airportCode, typeSensor, valueSensor, dateTime))
}

func writeFile(filePath string, txt string) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(txt + "\n"); err != nil {
		log.Println(err)
	}
}
