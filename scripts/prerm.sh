#!/bin/bash


if [ -f "/etc/systemd/system/hermes.service" ]; then
    systemctl stop hermes
    systemctl disable hermes
    systemctl daemon-reload
fi
