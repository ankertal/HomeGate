#!/bin/bash
if test -f "/tmp/gate-suspend.txt"; then
    exit
fi
if pgrep gate >/dev/null
then
  echo "Gate app is running."
else
  echo "Gate app is not running... Starting it now"
  cd /home/pi/HomeGate/
  git pull
  /home/pi/HomeGate/gate.py >> /tmp/gate.log  2>&1 &
fi