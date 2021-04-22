#!/usr/local/bin/python3
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


gateOperationLock = threading.Lock()

base_dir = "/home/pi/HomeGate"
signals_dir = base_dir + "/signals/"

deployment = "Tal"

OPEN_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
STOP_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
CLOSE_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]


NUM_ATTEMPTS = 5
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

def select_signal(action):
    return {
        Open: OPEN_TRANSMIT_SIGNAL,
        Close: STOP_TRANSMIT_SIGNAL,
        Close: CLOSE_TRANSMIT_SIGNAL
    }.get(action, OPEN_TRANSMIT_SIGNAL)


# Read a signal file into a signal value array
def read_signal(signal_file, signal):
    try:
        with open(signal_file, "r") as s:
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
    log_screen("Transmitting gate signal, length: " +
               str(len(signal[0])))
    for t in range(NUM_ATTEMPTS):
        for i in range(len(signal[0])):
            # print("going to transmit % d and to sleep %f secs" %
            #      (TRANSMIT_SIGNAL[1][i], TRANSMIT_SIGNAL[0][i]))
            GPIO.output(TRANSMIT_PIN, signal[1][i])
            time.sleep(signal[0][i])
    GPIO.cleanup()
    log_screen("Transmission done")


users = {'Tal': {}, 'Gilad': {}, 'Yaron': {}, 'Doron': {}}
users['Tal']['Tal'] = '024365645'
users['Tal']['Dorit'] = '028014405'
users['Yaron']['Tal'] = '024365645'
users['Gilad']['Gilad'] = '12345678'
users['Doron']['Doron'] = '12345678'

def main():
    # 1. read users from file and build dictionary {user:password}
    # 2. loop forever and read from db (see example)

    global gateOperationLock

    now = datetime.now()
    print('Starting HomeGate Service for Deployment: ' +
          deployment + " @ " + now.strftime("%d/%m/%Y %H:%M:%S"), flush=True)

    # start loop
    count = 1
    healthEntriesCounter = 0
    url = 'http://homegate.uaenorth.cloudapp.azure.com/status'
    myobj = {'deployment': 'Tal'}

    while True:
        gateOperationLock.acquire()
        x = requests.post(url, data = myobj)
        try:
            statusJson = json.load(x.text)
            button = statusJson['status']
            log_screen('Deployment: {0}'.format(deployment))
            log_screen('Button Press: {0}'.format(button))
            signal = select_signal(triggerAction)
            transmit_signal(signal)
        except:
            pass
        gateOperationLock.release()
        # Sleep a bit to avoid busy waiting
        time.sleep(0.5)


if __name__ == "__main__":
    open_signal_file = signals_dir + deployment + "-open.txt"
    stop_signal_file = signals_dir + deployment + "-stop.txt"
    close_signal_file = signals_dir + deployment + "-close.txt"

    read_signal(open_signal_file, OPEN_TRANSMIT_SIGNAL)
    read_signal(stop_signal_file, STOP_TRANSMIT_SIGNAL)
    read_signal(close_signal_file, CLOSE_TRANSMIT_SIGNAL)
    signal.signal(signal.SIGHUP, dump_state)

main()
