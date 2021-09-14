#!/usr/bin/env bash
[ "${FLOCKER}" != "$0" ] && exec env FLOCKER="$0" flock -en "$0" "$0" "$@" || :

do_log () {
   printf "%s\n" "$1 "
   echo $1  >>  /tmp/status
}

do_log "periodic attempting wifi-connect"

export DBUS_SYSTEM_BUS_ADDRESS=unix:path=/run/dbus/system_bus_socket
# Optional step - it takes couple of seconds (or longer) to establish a WiFi connection
# sometimes. In this case, following checks will fail and wifi-connect
# will be launched even if the device will be able to connect to a WiFi network.
# If this is your case, you can wait for a while and then check for the connection.
sleep 15

# Choose a condition for running WiFi Connect according to your use case:

# 1. Is there a default gateway?
# ip route | grep default

# 2. Is there Internet connectivity?
# nmcli -t g | grep full

# 3. Is there Internet connectivity via a google ping?
# wget --spider http://google.com 2>&1

# 4. Is there an active WiFi connection?
iwgetid -r
if [ $? -eq 0 ]; then
    do_log 'Skipping WiFi Connect'
else
    do_log 'Starting WiFi Connect'
    /usr/local/sbin/wifi-connect --portal-interface wlan0 --portal-ssid HomeGate >> /tmp/status
fi
do_log 'Done'
