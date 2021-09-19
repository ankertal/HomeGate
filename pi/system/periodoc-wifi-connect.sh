#!/bin/bash
touch /var/lock/wifi-connect.lock

for (( ; ; ))
do
   sudo /home/pi/HomeGate/system/start-wifi-connect.sh 
   sleep 60
done
