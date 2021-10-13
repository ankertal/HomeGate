# record and derive signal:
# We sample the GPIO and get a stream of values: <timestamp, val> where timestamp is the full time
# when the sample was taken and the value is 0 or 1
# Afterwards, we need to collect the samples into groups of consecutive values - aka segments
# i.e., if we have sample stream of [<t1, 0>, <t2, 0>, <t3, 1>, <t4, 1>, <t5, 0>...] we need to convert
# it to stream of <duration1=t3-t1, 0>, <duration2=t5-t3, 1> as the signal was 0 during the time
# interval of t3-t1 seconds and then it was 1 for the duration of t5-t3.
# once we have this segment stream we can continue and derive the minimal signal stream that will act as
# the RF remote control button:
# Every RF (433mhz) remote has a LONG silence in transmission between actually transmitting a button signal
# This long silence (probaby used for sync) is significantly larger then the regular signal duration.
# The difference is in order of magnitudes. We estimate the duration to be of 0.015 secs where a regular
# value transmitted usually takes 0.003-0.008 seconds.
# So we are looking for "silence" segments identified by duration > 0.01 sec, and mark them
# Later on we take two consecutive signals, starting with one "silence" signal, followed by the actual
# signal data and then takes the next signal including the post "silence" segment:
# the result will look like:  [LS][signal data][LS][signa data][LS]
# where LS is a long silence. We record this stream to a file, to be later used by a transmitter code
from datetime import datetime
from datetime import timedelta
import RPi.GPIO as GPIO

RECEIVED_SIGNAL = [[], []]  # [[time of reading], [signal reading deltaT]]
SIGNAL_SEGMENTS = [[], []]  # [[duration], [value 0/1]]
MAX_SILENCE_INDEXES = []
RECORD_DURATION = 1
RECEIVE_PIN = 24


def record_signal():
    print('1', flush=True)
    GPIO.setmode(GPIO.BCM)
    print('2', flush=True)
    GPIO.setup(RECEIVE_PIN, GPIO.IN)
    print('3', flush=True)

    cumulative_time = 0
    beginning_time = datetime.now()
    print('**Started recording**')

    while cumulative_time < RECORD_DURATION:
        sample_time = datetime.now()
        time_delta = sample_time - beginning_time
        val = GPIO.input(RECEIVE_PIN)
        RECEIVED_SIGNAL[0].append(time_delta)
        RECEIVED_SIGNAL[1].append(val)
        cumulative_time = time_delta.seconds
    print('**Ended recording**')
    print(len(RECEIVED_SIGNAL[0]), 'samples recorded')
    GPIO.cleanup()
    with open("/tmp/out_data.txt", "w") as f:
        for entry in range(len(RECEIVED_SIGNAL[0])):
            f.write("{0} {1}\n".format(
                RECEIVED_SIGNAL[0][entry].seconds +
                RECEIVED_SIGNAL[0][entry].microseconds/1000000.0, RECEIVED_SIGNAL[1][entry]))


def extract_signal_bursts():
    print('**Processing results**')
    prev_val = int(RECEIVED_SIGNAL[1][0])
    prev_sample = 0  # = RECEIVED_SIGNAL[0][0]
    # prev_timepoint = datetime.strptime(RECEIVED_SIGNAL[0][0], '%H:%M:%S.%f')
    with open("/tmp/out_segments.txt", "w") as f:
        for i in range(1, len(RECEIVED_SIGNAL[0])):
            # timepoint = datetime.strptime(RECEIVED_SIGNAL[0][i], '%H:%M:%S.%f')
            val = int(RECEIVED_SIGNAL[1][i])
            if prev_val != val:
                timepoint = RECEIVED_SIGNAL[0][i]
                # accumulated = timepoint - prev_timepoint
                # duration_formatted = accumulated.seconds + accumulated.microseconds/1000000.0
                duration = timepoint - RECEIVED_SIGNAL[0][prev_sample]
                SIGNAL_SEGMENTS[0].append(duration)
                SIGNAL_SEGMENTS[1].append(prev_val)
                f.write("{0} {1}\n".format(duration, prev_val))
                prev_val = val
                prev_sample = i
    print("segments length", len(SIGNAL_SEGMENTS[0]), "\n")
    if prev_sample != len(RECEIVED_SIGNAL[0]):
        f = open("/tmp/out_segments.txt", "a")
        duration = RECEIVED_SIGNAL[0][len(
            RECEIVED_SIGNAL[0]) - 1] - RECEIVED_SIGNAL[0][prev_sample]
        SIGNAL_SEGMENTS[0].append(duration)
        SIGNAL_SEGMENTS[1].append(prev_val)
        f.write("{0} {1}\n".format(duration, prev_val))
        f.close()


def extract_button_press(signal_file_name):
    LARGE_SILENCE_DURATION_MIN = 0.01
    for i in range(len(SIGNAL_SEGMENTS[0])):
        duration = SIGNAL_SEGMENTS[0][i].seconds + \
            SIGNAL_SEGMENTS[0][i].microseconds/1000000.0
        if duration > LARGE_SILENCE_DURATION_MIN:
            MAX_SILENCE_INDEXES.append(i)
    if len(MAX_SILENCE_INDEXES) < 3:
        print("Signal too short? could not find to large silence segments", flush=True)
    else:
        print("Number of large silence segments:", len(
            MAX_SILENCE_INDEXES), "\n", flush=True)
        print("going to take samples from ",
              MAX_SILENCE_INDEXES[0], "to ", MAX_SILENCE_INDEXES[2], "\n", flush=True)
        with open(signal_file_name, "w") as f:
            for i in range(MAX_SILENCE_INDEXES[0], MAX_SILENCE_INDEXES[2]):
                f.write("{0} {1}\n".format(
                    SIGNAL_SEGMENTS[0][i].seconds + SIGNAL_SEGMENTS[0][i].microseconds/1000000.0, SIGNAL_SEGMENTS[1][i]))

def learn_test():
    print("hello")

def record_button(signal_file_name):
    print("record_button", flush=True)
    record_signal()
    print("extracting burst", flush=True)
    extract_signal_bursts()
    print("extract button", flush=True)
    extract_button_press(signal_file_name)

if __name__ == '__main__':
    record_signal()
    extract_signal_bursts()
    extract_button_press(sys.argv[1])
