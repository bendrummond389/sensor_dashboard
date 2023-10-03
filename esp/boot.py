import network
import time
import json

def load_mqtt_config():
    try:
        with open('mqtt_config.json', 'r') as file:
            return json.load(file)
    except Exception as e:
        print(f"Could not read MQTT config: {e}")
        return {}

config = load_mqtt_config()

SSID = config.get("SSID", "default")
WIFI_PASSWORD  = config.get("WIFI_PASSWORD", "default")

def connect_wifi():
    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)
    ssid = SSID
    password = WIFI_PASSWORD

    retries = 5 
    
    while not wlan.isconnected() and retries > 0:
        print('Attempting to connect to network...')
        
        if not wlan.isconnected():
            wlan.connect(ssid, password) 

            for i in range(10):
                if wlan.isconnected():
                    break
                time.sleep(1)
        
        retries -= 1
        
        if wlan.isconnected():
            print(f'connected to {ssid}')
            print('Network config:', wlan.ifconfig())
            break 
        else:
            print(f'Failed to connect. Retries left: {retries}')

connect_wifi()
