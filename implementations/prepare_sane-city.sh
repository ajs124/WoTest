#!/usr/bin/env bash
set -e
cd "$(dirname "${BASH_SOURCE[0]}")"/sane-city
mvn package -DskipTests  # FIXME
