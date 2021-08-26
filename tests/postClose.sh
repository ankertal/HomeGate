#!/bin/bash

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"gate_name":"Tal","email":"talanker@gmail.com", "password":"024365645"}' \
  http://homegate.uaenorth.cloudapp.azure.com/close
