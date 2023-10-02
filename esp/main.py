import machine
import time

# Initialize pin 2 as an output pin
led = machine.Pin(2, machine.Pin.OUT)

while True:
    led.on()  # Turn LED on
    time.sleep(0.1)  # Wait for 1 second

    led.off()  # Turn LED off
    time.sleep(0.1)  # Wait for 1 second