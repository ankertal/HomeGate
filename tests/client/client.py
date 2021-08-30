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
    gate = {'gate_name': 'gate-46154121241', 'is_open': True}
    gate_json = json.dumps(gate)
    ws.send(gate_json)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        host = "ws://localhost/stream"
    else:
        host = sys.argv[1]
    headers = ['Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJlbWFpbCI6Ind5YXJvbkBnbWFpbC5jb20iLCJleHAiOjE2MzAzMjU0OTUsInJvbGUiOiJ1c2VyIn0.P2wBiUKLVklYJWF1p5GBezOsb1878R4jRqONUYVbZjU']
    ws = websocket.WebSocketApp(host,
                                on_message=on_message,
                                on_error=on_error,
                                on_close=on_close,
                                header=headers)
    ws.on_open = on_open
    ws.run_forever()
