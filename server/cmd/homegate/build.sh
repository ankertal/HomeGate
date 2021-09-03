#!/bin/bash

echo "stop homegate service"

sudo systemctl stop homegate
rm ./homegate

echo "clear homegate service journal"
sudo journalctl --rotate
sudo journalctl --vacuum-time=1s

echo "building..."
go build ./...


echo "starting homegate service"
sudo systemctl start homegate



