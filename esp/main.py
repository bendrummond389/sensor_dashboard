import time
import json
from sensor import read_sensor
from simple import MQTTClient
import machine
import ubinascii

HEARTBEAT = 0
HEARTBEAT_ACK = 1
SENSOR_DISCOVERY = 2
SENSOR_DISCOVERY_RESPONSE = 3
SENSOR_DATA = 4

client = None

def load_mqtt_config():
    try:
        with open('mqtt_config.json', 'r') as file:
            return json.load(file)
    except Exception as e:
        print(f"Could not read MQTT config: {e}")
        return {}

def create_payload(msg_type, device_id, data=None):
    return json.dumps({
        "type": msg_type,
        "device_id": device_id,
        "timestamp": str(time.time()),
        "data": data
    })

def handle_heartbeat(topic, msg):
    global client
    ack_payload = create_payload(HEARTBEAT_ACK, DEVICE_ID)
    client.publish(DISCOVERY_CHANNEL, ack_payload)

def mqtt_callback(topic, msg):
    global client
    try:
        message = json.loads(msg)
        if message.get("type") == HEARTBEAT:
            handle_heartbeat(topic, msg)
    except Exception as e:
        print(f"Failed to process message: {e}")

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
    global client
    client = connect_to_broker()
    if client:
        client.subscribe(DISCOVERY_CHANNEL)
        initial_payload = create_payload(SENSOR_DISCOVERY, DEVICE_ID, {"data_topic": SENSOR_DATA_TOPIC})
        client.publish(DISCOVERY_CHANNEL, initial_payload)
        
        while True:
            try:
                client.check_msg()
                sensor_value = read_sensor()
                sensor_data_payload = create_payload(SENSOR_DATA, DEVICE_ID, sensor_value)
                client.publish(SENSOR_DATA_TOPIC, sensor_data_payload)
                client.check_msg()
                time.sleep(10)
            except KeyboardInterrupt:
                break
            except:
                client = connect_to_broker()
                time.sleep(5)

if __name__ == "__main__":
    config = load_mqtt_config()
    BROKER_ADDRESS = config.get("BROKER_ADDRESS", "default_broker")
    BROKER_PORT = config.get("BROKER_PORT", 1883)
    SENSOR_DATA_TOPIC = config.get("SENSOR_DATA_TOPIC", "default_topic")
    DISCOVERY_CHANNEL = config.get("DISCOVERY_CHANNEL", "main")
    DEVICE_ID = ubinascii.hexlify(machine.unique_id())
    main()
