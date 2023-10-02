from machine import ADC, Pin

def read_sensor():
    try:
        print("About to read sensor...")
        adc = ADC(Pin(36))  # Connecting to GPIO36 (also known as ADC0)
        val = adc.read()
        print(f"Sensor read successfully: {val}")
        return val
    except Exception as e:
        print(f"Exception in read_sensor: {e}")
        return None