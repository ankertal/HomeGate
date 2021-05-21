#!/bin/bash
sudo systemctl stop homegate
pushd /home/wyaron/work/HomeGate/cmd/gate/homegate/
go build ./...
sudo journalctl --rotate && sudo journalctl --vacuum-time=1s

sudo systemctl start homegate
popd