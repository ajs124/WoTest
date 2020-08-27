#!/usr/bin/env bash
set -eux

cd "$(dirname "${BASH_SOURCE[0]}")"

# some tests require a running MQTT server
mosquitto -c mosquitto.conf &
trap "kill %1" EXIT

./prepare_node-wot.sh

./prepare_sane-city.sh

./prepare_wot-py.sh
