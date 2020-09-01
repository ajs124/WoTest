#!/usr/bin/env bash
set -ex
cd "$(dirname "${BASH_SOURCE[0]}")"/node-wot
npm install
npm run bootstrap
npm run build
cd ..
rm -rf node-wot-install
mkdir -p node-wot-install/node_modules/@node-wot
cd $_
find ../../../node-wot/packages -maxdepth 1 -mindepth 1 -type d -print0 | \
    while read -r -d "" package; do ln -s "$package" .; done
