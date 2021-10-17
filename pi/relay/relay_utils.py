import sys
import warnings
import time
from threading import Timer
import RPi.GPIO as GPIO

# our relay is controlled via GPIO pin 23 (BCM)
CTL_OUT = 23
    
def init_relay():       
    # switch to BCM
    GPIO.setmode(GPIO.BCM)

    # disable warning of used ports
    GPIO.setwarnings(False)

    # set the port as output port
    GPIO.setup(CTL_OUT, GPIO.OUT)

def start_relay():
    GPIO.output(CTL_OUT, GPIO.HIGH)

def stop_relay():
    GPIO.output(CTL_OUT, GPIO.LOW)

def cleanup():
    GPIO.cleanup()

def is_output_high():
    return GPIO.input(CTL_OUT)

def test_relay():
   """Test relay on and off cycle"""
   
   # check if the output is high
   print('current control output is: ' + str(is_output_high()) + ' (should be off)')

   # start the relay
   start_relay()

   print('current control output is: ' + str(is_output_high()) + ' (should be on)')
   
   # setup a timer to stop the relay after 5 seconds
   t = Timer(5, stop_relay)
   t.start()

   # wait for the timer to finish
   t.join()   

if __name__ == '__main__': 
    init_relay()
    test_relay()
    cleanup()