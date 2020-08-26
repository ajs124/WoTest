#!/usr/bin/env bash
set -eu
if [ -f /etc/NIXOS ]; then
    nix-shell --run "go run ."
elif [ $(which docker) ]; then
    docker build .
else
    requiredBinaries=("go" "node" "npm" "python" "pip")
    for b in $requiredBinaries; do
        which $b > /dev/null
        if [ $? != 0 ]; then
            echo "The required binary $b was not found in your PATH"
            exit 1
        fi
    done
    cd src
    go run .
fi
