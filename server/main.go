package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// DeviceInfo represents the JSON payload received from the sensor
type DeviceInfo struct {
	DeviceID  string `json:"device_id"`
	DataTopic string `json:"data_topic"`
}

// extractDeviceInfo parses the JSON payload to extract device information
func extractDeviceInfo(payload []byte) (*DeviceInfo, error) {
	var info DeviceInfo
	if err := json.Unmarshal(payload, &info); err != nil {
		return nil, err
	}
	return &info, nil
}


func main() {
	// Initialize a map to hold the last reading for each sensor
	var sensorData = make(map[string]string)
	var sensorDataMutex = &sync.Mutex{}

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
		// Extract device info from the JSON payload
		deviceInfo, err := extractDeviceInfo(msg.Payload())
		if err != nil {
			log.Printf("Failed to extract device info: %v", err)
			return
		}

		fmt.Printf("Received device info: %+v\n", deviceInfo)

		// Subscribe to the device's data topic
		client.Subscribe(deviceInfo.DataTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
			// Update the last reading for this sensor
			sensorDataMutex.Lock()
			sensorData[deviceInfo.DeviceID] = string(msg.Payload())
			sensorDataMutex.Unlock()

			fmt.Printf("Received data on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
			fmt.Printf("Current sensor data: %+v\n", sensorData)
		})
	})

	// Keep the application running
	select {}
}
