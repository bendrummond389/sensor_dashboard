import time
import json
from sensor import read_sensor
from simple import MQTTClient
import machine
import ubinascii

# Define message types
HEARTBEAT, HEARTBEAT_ACK, SENSOR_DISCOVERY, SENSOR_DISCOVERY_RESPONSE, SENSOR_DATA = range(5)

client = None  # Initialize MQTT client to None

def load_mqtt_config():
    """Load MQTT configuration from a file."""
    try:
        with open('mqtt_config.json', 'r') as file:
            return json.load(file)
    except Exception as e:
        print(f"Could not read MQTT config: {e}")
        return {}

def create_payload(msg_type, device_id, data=None):
    """Create a JSON payload for MQTT messages."""
    return json.dumps({
        "type": msg_type,
        "device_id": device_id,
        "timestamp": str(time.time()),
        "data": data
    })

def handle_heartbeat(topic, msg):
    """Handle heartbeat messages by sending an acknowledgment."""
    global client
    ack_payload = create_payload(HEARTBEAT_ACK, DEVICE_ID)
    client.publish(DISCOVERY_CHANNEL, ack_payload)

def mqtt_callback(topic, msg):
    """Callback function to handle incoming MQTT messages."""
    global client
    try:
        message = json.loads(msg)
        if message.get("type") == HEARTBEAT:
            handle_heartbeat(topic, msg)
    except Exception as e:
        print(f"Failed to process message: {e}")

def connect_to_broker():
    """Connect to the MQTT broker and return the client."""
    try:
        client = MQTTClient(DEVICE_ID, BROKER_ADDRESS, port=BROKER_PORT)
        client.set_callback(mqtt_callback)
        client.connect()
        return client
    except Exception as e:
        print(f"Exception during MQTT connection: {e}")
        return None

def main():
    """Main function to handle MQTT communication and sensor reading."""
    global client
    client = connect_to_broker()
    if client:
        client.subscribe(DISCOVERY_CHANNEL)
        initial_payload = create_payload(SENSOR_DISCOVERY, DEVICE_ID, {"data_topic": SENSOR_DATA_TOPIC})
        client.publish(DISCOVERY_CHANNEL, initial_payload)
        
        while True:
            try:
                # Check for new messages and publish sensor data
                client.check_msg()
                sensor_value = read_sensor()
                sensor_data_payload = create_payload(SENSOR_DATA, DEVICE_ID, sensor_value)
                client.publish(SENSOR_DATA_TOPIC, sensor_data_payload)
                client.check_msg()
                time.sleep(10)
            except KeyboardInterrupt:
                break  # Exit the loop on keyboard interrupt
            except:
                client = connect_to_broker()  # Reconnect to broker on any other exception
                time.sleep(5)

if __name__ == "__main__":
    config = load_mqtt_config()  # Load configuration
    # Set up global configuration variables
    BROKER_ADDRESS = config.get("BROKER_ADDRESS", "default_broker")
    BROKER_PORT = config.get("BROKER_PORT", 1883)
    SENSOR_DATA_TOPIC = config.get("SENSOR_DATA_TOPIC", "default_topic")
    DISCOVERY_CHANNEL = config.get("DISCOVERY_CHANNEL", "main")
    DEVICE_ID = ubinascii.hexlify(machine.unique_id())
    main()  # Run the main function
