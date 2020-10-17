package mqttredis

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/audoti/tp-infrastructure/pkg/putils"
	"github.com/gomodule/redigo/redis"
)

var redisconn redis.Conn
var mqttclient mqtt.Client

func redisAppend(key string, _time time.Time, value float64) {
	_, err := redisconn.Do("APPEND", key, fmt.Sprintf(",{\"v\":%f,\"d\":\"%s\"}", value, _time.Format(time.RFC3339)))
	if err != nil {
		log.Fatal(err)
	}
}

func redisAvg(key string, data float64) {
	temp, err := redis.Float64(redisconn.Do("GET", key))
	if err == nil {
		data = (data + temp) / 2
	}
	_, err = redisconn.Do("SET", key, data)
	if err != nil {
		log.Fatal(err)
	}
}

func redisMin(key string, data float64) {
	temp, err := redis.Float64(redisconn.Do("GET", key))
	if err == nil {
		data = math.Min(temp, data)
	}
	_, err = redisconn.Do("SET", key, data)
	if err != nil {
		log.Fatal(err)
	}
}

func redisMax(key string, data float64) {
	temp, err := redis.Float64(redisconn.Do("GET", key))
	if err == nil {
		data = math.Max(temp, data)
	}
	_, err = redisconn.Do("SET", key, data)
	if err != nil {
		log.Fatal(err)
	}
}

func redisIncr(key string) {
	_, err := redisconn.Do("INCR", key)
	if err != nil {
		log.Fatal(err)
	}
}

// ExtractMsgData incoming from IoT MQTT message
// sample topic: aeropot/wind
// sample message: 0000,AITA,Wind speed,5,2020-10-17T23:48:59+02:00
func ExtractMsgData(msg string) (string, float64, time.Time) {
	s := strings.Split(msg, ",")
	value, _ := strconv.ParseFloat(s[3], 64)
	t, _ := time.Parse(time.RFC3339, s[4])
	// aita, value, time
	return s[1], value, t
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	payload := message.Payload()
	topic := message.Topic()
	fmt.Printf("%s <-- %s\n", topic, payload)

	sensor := strings.Split(topic, "aeroport/")[1]
	aita, value, _t := ExtractMsgData(string(payload))
	date := putils.TimeToDate(_t)
	// Handle bad dates
	if date == "1-01-01" {
		return
	}

	sensorDataPrefix := "sensor"

	// --- Airport sensor stats
	fullkeyprefix := []string{
		// sensor|2020-09-30|CDG|wind -> sensor|2020-09-30|CDG|wind|avg <-- IoT data goes here (sensor|2020-09-30|CDG|wind|data)
		fmt.Sprintf("%s|%s|%s|%s", sensorDataPrefix, date, aita, sensor),
		// sensor|2020-09|CDG|wind -> sensor|2020-09|CDG|wind|avg
		fmt.Sprintf("%s|%s|%s|%s", sensorDataPrefix, date[:7], aita, sensor),
		// sensor|2020|CDG|wind -> sensor|2020|CDG|wind|avg
		fmt.Sprintf("%s|%s|%s|%s", sensorDataPrefix, date[:4], aita, sensor),
		// sensor|total|CDG|wind -> sensor|total|CDG|wind|avg
		fmt.Sprintf("%s|total|%s|%s", sensorDataPrefix, aita, sensor),
	}

	// --- Airport common stats (counts)
	airportkeyprefix := []string{
		// sensor|2020-09-30|CDG -> sensor|2020-09-30|CDG|count (only)
		fmt.Sprintf("%s|%s|%s", sensorDataPrefix, date, aita),
		// sensor|2020-09|CDG -> sensor|2020-09|CDG|count (only)
		fmt.Sprintf("%s|%s|%s", sensorDataPrefix, date[:7], aita),
		// sensor|2020|CDG -> sensor|2020|CDG|count (only)
		fmt.Sprintf("%s|%s|%s", sensorDataPrefix, date[:4], aita),
		// sensor|total|CDG -> sensor|total|CDG|count (only)
		fmt.Sprintf("%s|total|%s", sensorDataPrefix, aita),
	}

	// --- Global sensors stats
	globalkeysensorprefix := []string{
		// sensor|2020-09-30|wind -> sensor|2020-09-30|wind|avg
		fmt.Sprintf("%s|%s|%s", sensorDataPrefix, date, sensor),
		// sensor|2020-09|wind -> sensor|2020-09|wind|avg
		fmt.Sprintf("%s|%s|%s", sensorDataPrefix, date[:7], sensor),
		// sensor|2020|wind -> sensor|2020|wind|avg
		fmt.Sprintf("%s|%s|%s", sensorDataPrefix, date[:4], sensor),
		// sensor|total|wind -> sensor|total|wind|avg
		fmt.Sprintf("%s|total|%s", sensorDataPrefix, sensor),
	}

	// --- Global common stats (counts)
	globalkeyprefix := []string{
		// sensor|2020-09-30 -> sensor|2020-09-30|count (only)
		fmt.Sprintf("%s|%s", sensorDataPrefix, date),
		// sensor|2020-09 -> sensor|2020-09|count (only)
		fmt.Sprintf("%s|%s", sensorDataPrefix, date[:7]),
		// sensor|2020 -> sensor|2020|count (only)
		fmt.Sprintf("%s|%s", sensorDataPrefix, date[:4]),
		// sensor|total -> sensor|total|count (only)
		fmt.Sprintf("%s", sensorDataPrefix),
	}

	// Append IoT data
	redisAppend(fullkeyprefix[0]+"|data", _t, value)

	// Refresh average, minimum and maximum data
	for _, v := range fullkeyprefix {
		redisAvg(v+"|avg", value)
		redisMin(v+"|min", value)
		redisMax(v+"|max", value)
		redisIncr(v + "|count")
	}
	for _, v := range globalkeysensorprefix {
		redisAvg(v+"|avg", value)
		redisMin(v+"|min", value)
		redisMax(v+"|max", value)
		redisIncr(v + "|count")
	}

	for _, v := range airportkeyprefix {
		redisIncr(v + "|count")
	}
	for _, v := range globalkeyprefix {
		redisIncr(v + "|count")
	}
}

// SubscribeAndReact will subscribe an MQTT client to a topic and execute a function on incoming message
func SubscribeAndReact(mqttclient mqtt.Client, topic string, onMessageReceived func(client mqtt.Client, message mqtt.Message)) {
	if token := mqttclient.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("Connected to MQTT broker.")
		fmt.Println("Subscribed to the topic \"" + topic + "\". Reacting to incoming messages.")
	}
}

// RunMqttListenerRedis will run a MQTT client which will save incoming IoT messages to a Redis instance
func RunMqttListenerRedis() {
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	// mqtt.ERROR = log.New(os.Stdout, "", 0)

	// Connect to the Redis instance
	_redisconn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	redisconn = _redisconn

	// Close the redis connection on exit
	defer redisconn.Close()

	// Connect to the MQTT broker
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:18830")
	opts.SetClientID("listener-redis-1")
	opts.SetCleanSession(true)
	_mqttclient := mqtt.NewClient(opts)
	if token := _mqttclient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	mqttclient = _mqttclient

	// MQTT subscribe and react
	SubscribeAndReact(mqttclient, "aeroport/#", onMessageReceived)

	// Kill the process on SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
