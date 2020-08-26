#!/usr/bin/env bash
set -eux

cd "$(dirname "${BASH_SOURCE[0]}")"

# some tests require a running MQTT server
mosquitto -c mosquitto.conf &
trap "kill %1" EXIT

cd node-wot
npm install
npm run bootstrap
npm run build
cd ../node-wot-install
npm install

cd ../sane-city
ls -A
mvn install -DskipTests  # FIXME

cd ../wot-py
python3 setup.py build
#WOTPY_TESTS_MQTT_BROKER_URL=mqtt://127.0.0.1 python3 setup.py test
rm -rf install
mkdir install
pip install -t install .
