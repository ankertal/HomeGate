#!/bin/bash


curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"deployment":"Tal","user":"tal", "password":"024365645"}' \
  http://homegate.uaenorth.cloudapp.azure.com/open