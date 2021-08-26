import websocket
from threading import Thread
import time
import sys
import json


def on_message(ws, message):
    print(message)


def on_error(ws, error):
    print(error)


def on_close(ws, close_status_code, close_msg):
    print("### closed ###")


def on_open(ws):
    gate = {'gate_name': 'Tal', 'email': 'talanker@gmail.com',
            'password': '024365645'}
    gate_json = json.dumps(gate)
    ws.send(gate_json)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        host = "ws://localhost/stream"
    else:
        host = sys.argv[1]
    ws = websocket.WebSocketApp(host,
                                on_message=on_message,
                                on_error=on_error,
                                on_close=on_close)
    ws.on_open = on_open
    ws.run_forever()
