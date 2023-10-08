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

// MessageType enumerates the types of messages that can be sent/received.
type MessageType int

const (
	// Various message types
	Heartbeat               MessageType = iota // Heartbeat message type
	HeartbeatAck                               // Heartbeat acknowledgment message type
	SensorDiscovery                            // Sensor discovery message type
	SensorDiscoveryResponse                    // Sensor discovery response message type
	SensorData                                 // Sensor data message type
)

// Message struct represents the structure of messages exchanged.
type Message struct {
	Type      MessageType `json:"type"`      // The type of message
	DeviceID  string      `json:"device_id"` // The ID of the device sending/receiving the message
	Timestamp time.Time   `json:"timestamp"` // Timestamp when the message was sent/received
	Data      interface{} `json:"data"`      // The data contained in the message
}

// initMQTTClientWithRetry initializes an MQTT client with retry logic for connecting to the broker.
func initMQTTClientWithRetry(broker string, maxRetries int) mqtt.Client {
	var client mqtt.Client // MQTT client instance
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + broker).SetClientID("go-server")

	// Retry logic for connecting to the MQTT broker
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

	// Fatal log if the client fails to connect after max retries
	if client == nil || !client.IsConnected() {
		log.Fatal("Failed to connect to MQTT broker after max retries")
	}

	return client // Return the MQTT client
}

// sendHeartbeat sends a heartbeat message at regular intervals to the MQTT broker.
func sendHeartbeat(client mqtt.Client, mainTopic string, interval time.Duration) {
	ticker := time.NewTicker(interval) // Setup a ticker for the heartbeat interval
	for range ticker.C {
		heartbeatMessage := Message{
			Type:      Heartbeat,
			DeviceID:  "server",
			Timestamp: time.Now(),
			Data:      nil,
		}
		messageJSON, err := json.Marshal(heartbeatMessage) // Marshal the heartbeat message to JSON
		if err != nil {
			log.Printf("Failed to marshal heartbeat message: %v", err)
			continue
		}
		// Publish the heartbeat message to the MQTT broker
		token := client.Publish(mainTopic, 0, false, messageJSON)
		token.Wait()
		fmt.Println("Sent heartbeat pulse")
	}
}

// main is the entry point of the program.
func main() {
	mqttBroker := os.Getenv("MQTT_BROKER") // Get the MQTT broker address from the environment
	if mqttBroker == "" {
		mqttBroker = "localhost" // Default to localhost if no broker address is provided
	}
	if !strings.Contains(mqttBroker, ":") {
		mqttBroker = fmt.Sprintf("%s:1883", mqttBroker) // Append port number if not provided
	}

	client := initMQTTClientWithRetry(mqttBroker, 5) // Initialize the MQTT client
	fmt.Printf("Connected to MQTT broker at %s\n", mqttBroker)

	go sendHeartbeat(client, "main", 30*time.Second) // Start sending heartbeat messages

	select {} // Keep the main function running indefinitely
}
