import time
import machine
import json
from sensor import read_sensor
from simple import MQTTClient

def load_mqtt_config():
    try:
        with open('mqtt_config.json', 'r') as file:
            return json.load(file)
    except Exception as e:
        print(f"Could not read MQTT config: {e}")
        return {}

config = load_mqtt_config()

BROKER_ADDRESS = config.get("BROKER_ADDRESS", "default_broker")
BROKER_PORT = config.get("BROKER_PORT", 1883)
SENSOR_DATA_TOPIC = config.get("SENSOR_DATA_TOPIC", "default_topic")
DISCOVERY_CHANNEL = config.get("DISCOVERY_CHANNEL", "default")

DEVICE_ID = "water_sensor"

def mqtt_callback(topic, msg):
    print(f"Received message: {msg} on topic: {topic}")

def connect_to_broker():
    try:
        client = MQTTClient(DEVICE_ID, BROKER_ADDRESS, port=BROKER_PORT)
        client.set_callback(mqtt_callback)
        client.connect()
        return client
    except Exception as e:
        print(f"Exception during MQTT connection: {e}")
        return None

def main():
    client = connect_to_broker()
    if client:
        # Send initial sensor info
        initial_payload = {
            "device_id": DEVICE_ID,
            "data_topic": SENSOR_DATA_TOPIC,
        }
        client.publish(DISCOVERY_CHANNEL, json.dumps(initial_payload))
        
        while True:
            try:
                sensor_value = read_sensor()
                print(f"DEBUG: Water sensor value: {sensor_value}")
                client.publish(SENSOR_DATA_TOPIC, str(sensor_value))
                client.wait_msg()
                time.sleep(1)
            except KeyboardInterrupt:
                print("DEBUG: Disconnected from MQTT broker.")
                break
            except:
                print("DEBUG: An error occurred. Trying to reconnect.")
                client = connect_to_broker()
                time.sleep(5)

if __name__ == "__main__":
    main()
