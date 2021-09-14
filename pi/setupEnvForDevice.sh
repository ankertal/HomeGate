#!/bin/bash

if [[ "$#" -ne 2 ]]; then
    echo "this script needs exactly 2 parameters"
    exit 1
 fi

echo "email: $1";
echo "password: $2";

sed -i 's/GATE_USER=.*/GATE_USER='$1'/' /home/pi/work/HomeGate/pi/.env
sed -i 's/GATE_PASSWORD=.*/GATE_PASSWORD='$2'/' /home/pi/work/HomeGate/pi/.env

echo "killing gate to evaluate new environment..."
pkill -9 gate

