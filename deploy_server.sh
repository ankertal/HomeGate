#!/bin/bash
sudo systemctl stop homegate
pushd /home/ankertal/Work/HomeGate/cmd/gate/homegate/
go build ./...
sudo systemctl start homegate
popd