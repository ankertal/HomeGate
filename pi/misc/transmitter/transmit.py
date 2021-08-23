import time
import sys
import RPi.GPIO as GPIO

from datetime import timedelta
import RPi.GPIO as GPIO

TRANSMIT_SIGNAL = [[], []]  # [[length], [Value 0/1]]
NUM_ATTEMPTS = 20
TRANSMIT_PIN = 23


def read_signal(signal_file):
    with open(signal_file, "r") as s:
        for line in s:
            fields = line.split()
            signal_value = int(fields[1])
            signal_duration = float(fields[0])
            TRANSMIT_SIGNAL[0].append(signal_duration)
            TRANSMIT_SIGNAL[1].append(signal_value)


def transmit_signal():
    '''Transmit a signal stream using the GPIO transmitter'''
    GPIO.setmode(GPIO.BCM)
    GPIO.setup(TRANSMIT_PIN, GPIO.OUT)
    print("going to transmit % d times" % len(TRANSMIT_SIGNAL[0]))
    for t in range(NUM_ATTEMPTS):
        for i in range(len(TRANSMIT_SIGNAL[0])):
            # print("going to transmit % d and to sleep %f secs" %
            #      (TRANSMIT_SIGNAL[1][i], TRANSMIT_SIGNAL[0][i]))
            GPIO.output(TRANSMIT_PIN, TRANSMIT_SIGNAL[1][i])
            time.sleep(TRANSMIT_SIGNAL[0][i])
    GPIO.cleanup()


if __name__ == '__main__':
    read_signal(sys.argv[1])
    transmit_signal()
