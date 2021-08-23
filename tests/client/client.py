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
    deployment = {'deployment': 'Tal', 'user': 'tal',  'password': '024365645'}
    deployment_json = json.dumps(deployment)
    ws.send(deployment_json)


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
