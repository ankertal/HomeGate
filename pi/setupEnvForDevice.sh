#!/bin/bash

if [[ "$#" -ne 1 ]]; then
    echo "this script needs exactly 1 parameters (gate ID)"
    exit 1
 fi

echo "gate_id: $1";

sed -i 's/GATE_ID=.*/GATE_ID='$1'/' /home/pi/work/HomeGate/pi/.env

echo "killing gate to evaluate new environment..."
pkill -9 gate

