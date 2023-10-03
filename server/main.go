package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DeviceInfo struct {
	DeviceID  string `json:"device_id"`
	DataTopic string `json:"data_topic"`
}

var sensorData = make(map[string]string)
var sensorDataMutex = &sync.Mutex{}

func initMQTTClientWithRetry(broker string, maxRetries int) mqtt.Client {
	var client mqtt.Client
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + broker).SetClientID("go-server")

	for i := 0; i < maxRetries; i++ {
		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("Failed to connect to MQTT broker, attempt %d: %s", i+1, token.Error())
			time.Sleep(5 * time.Second)
		} else {
			log.Printf("Successfully connected to MQTT broker")
			break
		}
	}

	if client == nil || !client.IsConnected() {
		log.Fatal("Failed to connect to MQTT broker after max retries")
	}

	return client
}


func extractDeviceInfo(payload []byte) (*DeviceInfo, error) {
	var info DeviceInfo
	if err := json.Unmarshal(payload, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func handleSensorData(client mqtt.Client, msg mqtt.Message) {
	deviceInfo, err := extractDeviceInfo(msg.Payload())
	if err != nil {
		log.Printf("Failed to extract device info: %v", err)
		return
	}

	client.Subscribe(deviceInfo.DataTopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		sensorDataMutex.Lock()
		sensorData[deviceInfo.DeviceID] = string(msg.Payload())
		sensorDataMutex.Unlock()

		fmt.Printf("Received data on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
		fmt.Printf("Current sensor data: %+v\n", sensorData)
	})
}

func main() {
	mqttBroker := os.Getenv("MQTT_BROKER")
	if mqttBroker == "" {
		mqttBroker = "localhost"
	}
	if !strings.Contains(mqttBroker, ":") {
		mqttBroker = fmt.Sprintf("%s:1883", mqttBroker)
	}

	client := initMQTTClientWithRetry(mqttBroker, 5)
	fmt.Printf("Connected to MQTT broker at %s\n", mqttBroker)

	mainTopic := "sensor"
	client.Subscribe(mainTopic, 0, handleSensorData)

	select {}
}
