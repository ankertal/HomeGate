#!/bin/bash
if pgrep gate >/dev/null
then
  echo "Gate app is running."
else
  echo "Gate app is not running... Starting it now"
  cd /home/pi/work/HomeGate
  git pull
  /home/pi/work/HomeGate/gate.py >> /tmp/gate.log  2>&1 &
fi