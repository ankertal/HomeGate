#!/usr/bin/python3

import sys
import os
import shutil
from datetime import datetime
from datetime import timedelta
import RPi.GPIO as GPIO
import threading
import json
import time
import faulthandler
import signal
import requests
import websocket
from dotenv import load_dotenv
from learn import learn_utils
from relay import relay_utils

# use dot env for simply development and testing
load_dotenv()

SIGNAL_DIR = os.getenv('SIGNAL_DIR')
SERVER_HOMEGATE_URL = os.getenv('SERVER_HOMEGATE_URL')
SIGNAL_PREFIX = os.getenv('SIGNAL_PREFIX')
GATE_ID = os.getenv('GATE_ID')
WIRED_MODE = os.getenv('WIRED_MODE')

if SIGNAL_PREFIX == None or SIGNAL_PREFIX == "":
    SIGNAL_PREFIX = 'default'

OPEN_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
STOP_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
CLOSE_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
CANDIDATE_OPEN_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
CANDIDATE_STOP_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
CANDIDATE_CLOSE_TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]


def dump_state(signalNumber, frame):
    faulthandler.dump_traceback(file=sys.stderr, all_threads=True)


def logit(*args, **kwargs):
    args = (datetime.now().strftime("%d/%m/%Y %H:%M:%S"), ) + args
    print(*args, flush=True)


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
        logit("warning", 'signal file:', signal_file, 'does not exist')


def transmit_signal(signal):
    NUM_ATTEMPTS = 7
    TRANSMIT_PIN = 23

    GPIO.setmode(GPIO.BCM)
    GPIO.setup(TRANSMIT_PIN, GPIO.OUT)

    if WIRED_MODE=='1':
        # we use a relay, simply close it
        logit("Using relay mode, closing relay for 1 sec...")
        relay_utils.start_relay()
        time.sleep(1)
        logit("Using relay mode, open relay for 1 sec...")
        relay_utils.stop_relay()
    else:
        logit("Transmitting gate signal, length: " + str(len(signal[0])))
        for t in range(NUM_ATTEMPTS):
            for i in range(len(signal[0])):
                GPIO.output(TRANSMIT_PIN, signal[1][i])
                time.sleep(signal[0][i])

    GPIO.cleanup()
    logit("Transmitting gate command [OK]")


def learn_open():
    now = datetime.now()
    logit('learn to open pressed: ' + SIGNAL_PREFIX + " @ " +
          now.strftime("%d/%m/%Y %H:%M:%S"))
    signal_file_name = "/tmp/" + SIGNAL_PREFIX + "-open.txt"

    logit("before record")
    logit("file name " + signal_file_name)
    learn_utils.record_button(signal_file_name)


def learn_close():
    now = datetime.now()
    logit('learn to close pressed: ' + SIGNAL_PREFIX + " @ " +
          now.strftime("%d/%m/%Y %H:%M:%S"))

    signal_file_name = "/tmp/" + SIGNAL_PREFIX + "-close.txt"
    learn_utils.record_button(signal_file_name + "-close.txt")
    read_signal(signal_file_name, CANDIDATE_CLOSE_TRANSMIT_SIGNAL)


def test_open():
    signal_file_name = "/tmp/" + SIGNAL_PREFIX + "-open.txt"
    read_signal(signal_file_name, CANDIDATE_OPEN_TRANSMIT_SIGNAL)
    transmit_signal(CANDIDATE_OPEN_TRANSMIT_SIGNAL)


def test_close():
    signal_file_name = "/tmp/" + SIGNAL_PREFIX + "-close.txt"
    read_signal(signal_file_name, CANDIDATE_CLOSE_TRANSMIT_SIGNAL)
    transmit_signal(CANDIDATE_CLOSE_TRANSMIT_SIGNAL)


def set_open():
    # copy org to backup, copy candidate to main
    now = datetime.now()
    dst_open_signal_file = SIGNAL_DIR + SIGNAL_PREFIX + "-open.txt"
    src = "/tmp/" + SIGNAL_PREFIX + "-open.txt"
    backup_open_signal_file = SIGNAL_DIR + \
        now.strftime("%d-%m-%Y-%H-%M-%S-") + SIGNAL_PREFIX + "-open.txt"
    try:
        shutil.copyfile(dst_open_signal_file, backup_open_signal_file)
    except:
        logit(
            'Could not backup original open gate signal - probably not recorded yet'
        )
        pass
    try:
        shutil.copyfile(src, dst_open_signal_file)
        read_signal(dst_open_signal_file, OPEN_TRANSMIT_SIGNAL)
    except:
        logit('Could not set open gate signal - probably not recorded yet')
        pass


def set_close():
    # copy org to backup, copy candidate to main
    now = datetime.now()
    dst_close_signal_file = SIGNAL_DIR + SIGNAL_PREFIX + "-close.txt"
    src = "/tmp/" + SIGNAL_PREFIX + "-close.txt"
    backup_close_signal_file = SIGNAL_DIR + \
        now.strftime("%d-%m-%Y-%H-%M-%S-") + SIGNAL_PREFIX + "-close.txt"
    try:
        shutil.copyfile(dst_close_signal_file, backup_close_signal_file)
    except:
        logit(
            'Could not backup original close gate signal - probably not recorded yet'
        )
        pass
    try:
        shutil.copyfile(src, dst_close_signal_file)
        read_signal(dst_close_signal_file, CLOSE_TRANSMIT_SIGNAL)
    except:
        logit('Could not set close gate signal - probably not recorded yet')
        pass


def on_remote_control(ws, button):
    logit("received a remote controler event: ", button)
    if button == "Close":
        transmit_signal(CLOSE_TRANSMIT_SIGNAL)
    if button == "Open":
        transmit_signal(OPEN_TRANSMIT_SIGNAL)
    if button == "LearnOpen":
        learn_open()
    if button == "LearnClose":
        learn_close()
    if button == "SetOpen":
        set_open()
    if button == "SetClose":
        set_close()
    if button == "TestOpen":
        test_open()
    if button == "TestClose":
        test_close()


def on_error(ws, error):
    logit("received an error from the homegate server", str(error))
    sys.exit("exiting and let cron restart")


def on_close(ws, close_status_code, close_msg):
    logit("### websocket to server closed ###")
    sys.exit("exiting and let cron restart")


def on_open(ws):
    postData = {}
    postData['gate_id'] = GATE_ID
    open_post_data_json = json.dumps(postData)
    ws.send(open_post_data_json)


def on_ping(wsapp, message):
    logit(
        "Got a ping from server, a pong reply has already been automatically sent."
    )


def load_ctrl_signals():
    open_signal_file = SIGNAL_DIR + SIGNAL_PREFIX + "-open.txt"
    stop_signal_file = SIGNAL_DIR + SIGNAL_PREFIX + "-stop.txt"
    close_signal_file = SIGNAL_DIR + SIGNAL_PREFIX + "-close.txt"

    read_signal(open_signal_file, OPEN_TRANSMIT_SIGNAL)
    read_signal(stop_signal_file, STOP_TRANSMIT_SIGNAL)
    read_signal(close_signal_file, CLOSE_TRANSMIT_SIGNAL)


if __name__ == "__main__":
    signal.signal(signal.SIGHUP, dump_state)

    load_ctrl_signals()
    now = datetime.now()
    streamEndpoint = "ws://" + SERVER_HOMEGATE_URL + "/stream"

    logit('Starting HomeGate Device: {0}, server: {1}, time: {2}'.format(
        SIGNAL_PREFIX, now.strftime("%d/%m/%Y %H:%M:%S"), streamEndpoint))
    ws = websocket.WebSocketApp(streamEndpoint,
                                on_message=on_remote_control,
                                on_ping=on_ping,
                                on_error=on_error,
                                on_close=on_close)
    ws.on_open = on_open
    ws.run_forever()
