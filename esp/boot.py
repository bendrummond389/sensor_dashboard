import network
import time

def connect_wifi():
    wlan = network.WLAN(network.STA_IF)
    wlan.active(True)

    retries = 5 
    
    while not wlan.isconnected() and retries > 0:
        print('Attempting to connect to network...')
        
        if not wlan.isconnected():
            wlan.connect('CasaLindaFiber', '10Strings') 

            for i in range(10):
                if wlan.isconnected():
                    break
                time.sleep(1)
        
        retries -= 1
        
        if wlan.isconnected():
            print('Network config:', wlan.ifconfig())
            break 
        else:
            print(f'Failed to connect. Retries left: {retries}')

connect_wifi()
