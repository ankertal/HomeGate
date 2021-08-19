#!/bin/bash

oldstamp=$(date +%s -d '1 day ago')

while IFS=\: read -r contype timestamp uuid
do
  echo $contype $uuid
  case "$contype" in
    802-11-wireless)
      if ((timestamp < oldstamp)); then
        nmcli con delete uuid "$uuid"
      else
        echo "$uuid: too new: skipping"
      fi
    ;;

    *)
      echo "skipping connection type $contype"
    ;;
  esac
done < <(nmcli -t -f TYPE,TIMESTAMP,UUID con)
