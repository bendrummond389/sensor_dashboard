# IoT Home Automation Messaging System Documentation

## Overview
This document describes the standardized messaging system used for communication between the Go server and the MicroPython-based microcontrollers in the IoT Home Automation project.

## Message Schema
All messages should adhere to the following JSON schema:

```json
{
  "type": "message_type",
  "device_id": "unique_device_id",
  "timestamp": "ISO_timestamp",
  "data": {
    // Type-specific data
  }
}
```

### Fields
- type: String that indicates the type of the message.
- device_id: Unique identifier for the device sending or receiving the message.
- timestamp: ISO 8601 formatted timestamp indicating when the message was sent.
- data: An object containing additional data specific to the message type.

## Message Types

### Heartbeat

- Type: heartbeat
- Direction: Server to Microcontroller
- Data Fields: None
- Description: Sent by the server at regular intervals to check device availability.

### Heartbeat Acknowledgment

- Type: heartbeat_ack
- Direction: Microcontroller to Server
- Data Fields: None
- Description: Sent by the microcontroller in response to a heartbeat.

### Sensor Discovery
- Type: sensor_discovery
- Direction: Server to Microcontroller
- Data Fields: None
- Description: Sent by the server on startup to discover active sensors.

### Sensor Discovery Response
- Type: sensor_discovery_response
- Direction: Microcontroller to Server
- Data Fields:
- data_topic: The MQTT topic the sensor will publish data to.
- Description: Sent by the microcontroller in response to a sensor discovery request.

### Sensor Data
- Type: sensor_data
- Direction: Microcontroller to Server
- Data Fields:
- temperature: Temperature reading.
- humidity: Humidity reading.
- Description: Regular sensor data sent by the microcontroller.

## Example Messages

### Heartbeat
```json
{
  "type": "heartbeat",
  "device_id": "server",
  "timestamp": "2023-10-03T10:15:30Z",
  "data": {}
}
```

### Sensor Discovery Response
```json
{
  "type": "sensor_discovery_response",
  "device_id": "device_123",
  "timestamp": "2023-10-03T10:16:30Z",
  "data": {
    "data_topic": "sensor/device_123/data"
  }
}
```