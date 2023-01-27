#!/bin/bash
set -euo pipefail
IFS=$'\n\t'
logfile="$1"

if [[ ! -f "$logfile" ]]; then
  echo "File $logfile does not exist"
  exit 1
fi

while true; do
  sleep $(shuf -i 1-3 -n 1);
  echo "connected $((1 + RANDOM % 1000))" | tee -a $logfile
done