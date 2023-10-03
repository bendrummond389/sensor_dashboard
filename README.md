# IoT Home Automation Project

## Overview

This project aims to simplify home automation by providing an easy-to-use IoT solution. It utilizes microcontrollers running MicroPython, MQTT for data transmission, and a Next.js dashboard for monitoring and control. The entire stack can be containerized using Docker, making it easy to deploy even on a Raspberry Pi.

## Features

- Microcontrollers running MicroPython
- MQTT Broker for real-time communication
- WebSocket server for real-time data availability
- Next.js Dashboard for monitoring and control
- Dockerized for easy deployment

## Prerequisites

- MicroPython-compatible microcontroller (e.g., ESP32, ESP8266)
- Docker
- Node.js and npm (for Next.js app)

## Quick Start

1. **Clone the Repository**

    ```bash
    git clone https://github.com/bendrummond389/sensor_dashboard
    ```

2. **Start Docker Containers**
    Navigate to the project directory and run:

    ```bash
    docker-compose up -d
    ```

3. **Configure Microcontroller**

    Modify the `config.json` file on the microcontroller with your MQTT broker URL, SSID, and password.
    Add your sensor specific code to sensor.py.

4. **Access Dashboard**

    Open your browser and navigate to the Next.js dashboard URL (usually `http://localhost:3000`).

## Configuration

### Microcontroller

Edit the `config.json` file to include:

- `mqtt_broker_url`: URL of your MQTT broker
- `ssid`: Your WiFi SSID
- `password`: Your WiFi password

Additionally, modify the `sensor.py` file to fit your specific sensor setup.

### MQTT Broker

Ensure your MQTT broker is running and accessible from the microcontroller and your network.

### Next.js Dashboard

No additional configuration needed. Just ensure it's running and accessible.

## Contributing

Feel free to open issues or submit pull requests. Your contributions are welcome!

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

