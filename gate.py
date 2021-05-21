#!/usr/bin/python3
from datetime import datetime
from datetime import timedelta
import RPi.GPIO as GPIO
import sys
import os
import threading
import json
import time
import threading
import faulthandler
import signal
import requests


base_dir = "/home/pi/work/HomeGate"
signals_dir = base_dir + "/signals/"

deployment = "Yaron"

OPEN_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
STOP_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
CLOSE_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]


NUM_ATTEMPTS = 20
TRANSMIT_PIN = 23


def dump_state(signalNumber, frame):
    faulthandler.dump_traceback(file=sys.stderr, all_threads=True)


def log_screen(*args, **kwargs):
    args = (datetime.now().strftime("%d/%m/%Y %H:%M:%S"),)+args
    print(*args, flush=True)

# In principle we currently support 3 button types - open/close/stop. This will translate into three
# signals to transmit (depending on the required action)
# select_signal function takes an action string and return the correspinding signal to transmit. If no
# match was found, then the default chosen signal is for the "open" action

Unknown = '0'
Close = '1'
Open = '2'
Stop = '3'
Update = '4'

def select_signal(action):
    return {
        Open: OPEN_TRANSMIT_SIGNAL,
        Stop: STOP_TRANSMIT_SIGNAL,
        Close: CLOSE_TRANSMIT_SIGNAL
    }.get(action, OPEN_TRANSMIT_SIGNAL)


# Read a signal file into a signal value array
def read_signal(signal_file, signal):
    try:
        with open(signal_file, "r") as s:
            print('===> signal file: ' + signal_file)
            for line in s:
                fields = line.split()
                signal_value = int(fields[1])
                signal_duration = float(fields[0])
                signal[0].append(signal_duration)
                signal[1].append(signal_value)

    except IOError:
        print('signal file: ' + signal_file + ' does not exist', flush=True)


# transmit_signal: Transmit a signal (can be signal to open/stop/close)
def transmit_signal(signal):
    global TRANSMIT_PIN
    global NUM_ATTEMPTS

    '''Transmit a signal stream using the GPIO transmitter'''
    GPIO.setmode(GPIO.BCM)
    GPIO.setup(TRANSMIT_PIN, GPIO.OUT)
    log_screen("Transmitting gate signal, length: " + str(len(signal[0])) + " times = " + str(NUM_ATTEMPTS))
    for t in range(NUM_ATTEMPTS):
        for i in range(len(signal[0])):
            # print("going to transmit % d and to sleep %f secs" %
            #      (TRANSMIT_SIGNAL[1][i], TRANSMIT_SIGNAL[0][i]))
            GPIO.output(TRANSMIT_PIN, signal[1][i])
            time.sleep(signal[0][i])
        time.sleep(0.5)
    GPIO.cleanup()
    log_screen("Transmission done")


def main():
    # 1. read users from file and build dictionary {user:password}
    # 2. loop forever and read from db (see example)

    now = datetime.now()
    print('Starting HomeGate Service for Deployment: ' +
          deployment + " @ " + now.strftime("%d/%m/%Y %H:%M:%S"), flush=True)

    # start loop
    healthEntriesCounter = 0
    url = 'http://weinsgate.uaenorth.cloudapp.azure.com/status'
    getGateRCStatus = '{"deployment": "Yaron", "user": "gate", "password": "homegate;415231"}'

    x = requests.get('http://weinsgate.uaenorth.cloudapp.azure.com/')
    log_screen('{0}'.format(x.text))

    while True:
         # Sleep a bit to avoid busy waiting
        time.sleep(2)
        x = requests.post(url, data = getGateRCStatus)
        try:
            statusJson = json.loads(x.text)
            button = statusJson['status']
            log_screen('Deployment: {0}'.format(deployment))
            log_screen('Button Press: {0}'.format(button))
            if button == Unknown:
                continue
            if button == Update: 
                return
            signal = select_signal(button)
            print(signal)
            transmit_signal(signal)
        except:
            print('exception\n', flush=True)

       


if __name__ == "__main__":
    open_signal_file = signals_dir + deployment + "-open.txt"
    stop_signal_file = signals_dir + deployment + "-stop.txt"
    close_signal_file = signals_dir + deployment + "-close.txt"

    read_signal(open_signal_file, OPEN_TRANSMIT_SIGNAL)
    read_signal(stop_signal_file, STOP_TRANSMIT_SIGNAL)
    read_signal(close_signal_file, CLOSE_TRANSMIT_SIGNAL)
    signal.signal(signal.SIGHUP, dump_state)

main()
