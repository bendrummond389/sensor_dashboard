import time
import machine
import json
from sensor import read_sensor
from simple import MQTTClient

def load_config():
    try:
        with open('mqtt_config.json', 'r') as file:
            return json.load(file)
    except Exception as e:
        print(f"Could not read MQTT config: {e}")
        return {}

config = load_config()

MQTT_BROKER = config.get("MQTT_BROKER", "default_broker")
MQTT_PORT = config.get("MQTT_PORT", 1883)
MQTT_TOPIC = config.get("MQTT_TOPIC", "default_topic")

CLIENT_ID = "water_sensor"

def mqtt_callback(topic, msg):
    print(f"Received message: {msg} on topic: {topic}")

def connect_mqtt():
    try:
        client = MQTTClient(CLIENT_ID, MQTT_BROKER, port=MQTT_PORT)
        client.set_callback(mqtt_callback)
        client.connect()
        client.subscribe(MQTT_TOPIC)
        return client
    except Exception as e:
        print(f"Exception during MQTT connection: {e}")
        return None

def main():
    client = connect_mqtt()
    if client:
        while True:
            try:
                sensor_value = read_sensor()
                print(f"DEBUG: Water sensor value: {sensor_value}")
                client.publish(MQTT_TOPIC, str(sensor_value))
                client.wait_msg()
                time.sleep(1)
            except KeyboardInterrupt:
                print("DEBUG: Disconnected from MQTT broker.")
                break
            except:
                print("DEBUG: An error occurred. Trying to reconnect.")
                client = connect_mqtt()
                time.sleep(5)

if __name__ == "__main__":
    main()
