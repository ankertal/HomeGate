#!/bin/bash
printf "Resetting device to factory defaults, removing all ssids\n"

printf "First, stop the wifi-connect service so we wont interfere with it\n"
systemctl stop wifi-connect
printf "wifi-connect is stopped\n"

printf "Now removing existing ssids\n"
while IFS=\: read -r contype timestamp uuid
do
  echo $contype $uuid
  case "$contype" in
    802-11-wireless)
      nmcli con delete uuid "$uuid"
    ;;

    *)
    ;;
  esac
done < <(nmcli -t -f TYPE,TIMESTAMP,UUID con)

printf "Done removing all ssids\n"

printf "Now, start again the wifi-connect service\n"

systemctl start wifi-connect

printf "Done\n"

