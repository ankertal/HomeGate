#!/bin/bash
if test -f "/tmp/gate-suspend.txt"; then
    exit
fi
if pgrep gate >/dev/null
then
  echo "Gate app is running."
else
  echo "Gate app is not running... Starting it now"
  cd /home/pi/work/HomeGate
  git pull
  /home/pi/work/HomeGate/pi/gate.py >> /tmp/gate.log  2>&1 &
  

  pkill -9 flask
  export FLASK_APP=/home/pi/work/HomeGate/pi/app.py
  echo "run homegate configuration server..."
  flask run --host=0.0.0.0 >> /tmp/gate-config.log  2>&1 &
fi