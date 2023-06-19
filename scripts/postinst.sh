#!/bin/bash


if [ -f "/etc/systemd/system/hermes.service" ]; then
    systemctl start hermes
    systemctl enable hermes
    systemctl daemon-reload
fi
