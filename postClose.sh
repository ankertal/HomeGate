#!/bin/bash


curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"deployment":"Yaron","user":"Yaron", "password":"029607017"}' \
  http://weinsgate.uaenorth.cloudapp.azure.com/close

  