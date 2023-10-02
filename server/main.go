package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Initialize MQTT broker settings
	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		mqttBroker = "localhost" // Default broker address
	}

	// Append port if not specified
	if !strings.Contains(mqttBroker, ":") {
		mqttBroker = fmt.Sprintf("%s:1883", mqttBroker)
	}
	fmt.Printf("MQTT Broker set to %s\n", mqttBroker)

	// Set up MQTT client options
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + mqttBroker).SetClientID("go-server")

	// Connect to MQTT broker
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Failed to connect to MQTT broker:", token.Error())
	}
	fmt.Printf("Connected to MQTT broker at %s\n", mqttBroker)

	// Subscribe to the main channel
	mainTopic := "sensor"
	client.Subscribe(mainTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Print the received message
		fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
	})

	// Keep the application running
	select {}
}
