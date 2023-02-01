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
  echo "Feb 01 15:51:45 light dnsmasq-dhcp[123456]: DHCPACK(wlp0s20f3) 10.42.0.227 aa:bb:cc:dd:ee:ff Pixel-5" | tee -a $logfile
done