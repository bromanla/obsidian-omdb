#!/bin/bash
# Common settings for scripts

SERVICE_NAME="omdb-bot"
UNIT_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

# Colored logging
success() { echo -e "\033[32m[OK]    $*\033[0m"; }   # Green
info()    { echo -e "\033[34m[INFO]  $*\033[0m"; }   # Blue
warn()    { echo -e "\033[33m[WARN]  $*\033[0m"; }   # Yellow
error()   { echo -e "\033[31m[ERROR] $*\033[0m" >&2; } # Red