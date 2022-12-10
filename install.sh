#!/bin/bash
if (systemctl -q is-active getgo.service)
    then
    systemctl stop getgo
fi
mkdir -p /apps/getgo/bin/
cp -u ./bin/getgo /apps/getgo/bin/getgo
cp getgo.service /etc/systemd/system/getgo.service
systemctl daemon-reload
systemctl start getgo