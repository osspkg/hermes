#!/bin/bash


if ! [ -d /var/lib/hermes/ ]; then
    mkdir /var/lib/hermes
fi

if [ -f "/etc/systemd/system/hermes.service" ]; then
    systemctl stop hermes
    systemctl disable hermes
    systemctl daemon-reload
fi
