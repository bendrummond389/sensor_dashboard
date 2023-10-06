package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MessageType int

const (
	Heartbeat MessageType = iota
	HeartbeatAck
	SensorDiscovery
	SensorDiscoveryResponse
	SensorData
)

type Message struct {
	Type      MessageType `json:"type"`
	DeviceID  string      `json:"device_id"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

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

func sendHeartbeat(client mqtt.Client, mainTopic string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		heartbeatMessage := Message{
			Type:      Heartbeat,
			DeviceID:  "server",
			Timestamp: time.Now(),
			Data:      nil,
		}
		messageJSON, err := json.Marshal(heartbeatMessage)
		if err != nil {
			log.Printf("Failed to marshal heartbeat message: %v", err)
			continue
		}
		token := client.Publish(mainTopic, 0, false, messageJSON)
		token.Wait()
		fmt.Println("Sent heartbeat pulse")
	}
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

	go sendHeartbeat(client, "main", 30*time.Second)

	select {}
}
