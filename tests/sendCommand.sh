#!/bin/bash

############################################################
# Help                                                     #
############################################################
Help() {
  echo "sendCommand make it easy to send gate command for testing"
  echo
  echo "Syntax: "
  echo "sendCommand -c [is_open|is_close|is_learn_open|is_learn_close|is_test_open|is_test_close|is_set_open|is_set_close]"
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

TOKEN=$(curl -s -X POST -H 'Accept: application/json' -H 'Content-Type: application/json' \
  --data '{"email":"wyaron@gmail.com","password":"123456"}' http://localhost/signin | jq -r '.accessToken')

GATE_NAME=$(curl -s -X POST -H 'Accept: application/json' -H 'Content-Type: application/json' \
  --data '{"email":"wyaron@gmail.com","password":"123456"}' http://localhost/signin | jq -r '.my_gate')

echo "Got JWT token: $TOKEN"
echo "Gate Name: $GATE_NAME"
echo "Command Is: $Command"

generate_post_data() {
  cat <<EOF
{
  "gate_name":"$GATE_NAME", 
  "$Command":true
}
EOF
}

curl -X POST -H 'Accept: application/json' -H "Authorization: Bearer ${TOKEN}" -H "Token: ${TOKEN}" -H 'Content-Type: application/json' \
  --data "$(generate_post_data)" http://localhost/command
