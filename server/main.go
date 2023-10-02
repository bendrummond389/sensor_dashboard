package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// type SensorInfo struct {
// 	Device    string `json:"device"`
// 	MqttTopic string `json:"mqtt_topic"`
// }

func main() {
	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		mqttBroker = "localhost" // Default
	}

	// append port if not already specified
	if !strings.Contains(mqttBroker, ":") {
		mqttBroker = fmt.Sprintf("%s:1883", mqttBroker)
	}
	fmt.Printf("MQTT Broker set to %s \n", mqttBroker)


	// set up client 
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + mqttBroker).SetClientID("go-server")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	fmt.Printf("Connected to MQTT broker at %s\n", mqttBroker)


	// connect to main channel 
	mainTopic := "sensor"
	client.Subscribe(mainTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Print the raw JSON payload
		fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
	})


	select{}

}