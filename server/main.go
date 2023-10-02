package main

import (
	"fmt"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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
		panic(token.Error())
	}

	fmt.Printf("Connected to MQTT broker at %s\n", mqttBroker)


	// set up listener on sensor topic and listen for device connections







}