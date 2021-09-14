#!/bin/bash

############################################################
# Help                                                     #
############################################################
Help() {
  echo "sendSiriCommand check the SR siri interface"
  echo
  echo "Syntax: "
  echo "sendSiriCommand -c [open|close]"
  echo
}

PARAMS=""
while (("$#")); do
  case "$1" in
  -h | --help)
    Help
    exit
    ;;
  -c | --command)
    if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
      Command=$2
      shift 2
    else
      echo "Error: Argument for $1 is missing" >&2
      exit 1
    fi
    ;;
  -* | --*=) # unsupported flags
    echo "Error: Unsupported flag $1" >&2
    exit 1
    ;;
  *) # preserve positional arguments
    PARAMS="$PARAMS $1"
    shift
    ;;
  esac
done

# set positional arguments in their proper place
eval set -- "$PARAMS"

generate_post_data() {
  cat <<EOF
{
  "gate_id":"gate-1245406299518", 
  "gate_command" : "$Command"
}
EOF
}

curl -X POST -H 'Accept: application/json' -H 'Content-Type: application/json' \
  --data "$(generate_post_data)" http://localhost/siri
