#!/bin/bash
touch /var/lock/wifi-connect.lock

for (( ; ; ))
do
   sudo /home/pi/HomeGate/pi/system/start-wifi-connect.sh 
   sleep 60
done
