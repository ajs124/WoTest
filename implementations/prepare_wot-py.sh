#!/usr/bin/env bash
set -ex
cd "$(dirname "${BASH_SOURCE[0]}")"/wot-py
python3 setup.py build
#WOTPY_TESTS_MQTT_BROKER_URL=mqtt://127.0.0.1 python3 setup.py test
rm -rf install
mkdir install
pip install -t install .
