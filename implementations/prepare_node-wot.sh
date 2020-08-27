#!/usr/bin/env bash
set -ex
cd "$(dirname "${BASH_SOURCE[0]}")"/node-wot
npm install
npm run bootstrap
npm run build
cd ../node-wot-install
npm install
