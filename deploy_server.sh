#!/bin/bash
sudo systemctl stop homegate
pushd /home/wyaron/work/HomeGate/cmd/gate/homegate/
go build ./...
sudo systemctl start homegate
popd