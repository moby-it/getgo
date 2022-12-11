#!/bin/bash

# check if docker exists in system
if ! command -v docker &> /dev/null
then
    echo "docker not found"
    exit 1
fi
echo "Installing..."
if (systemctl -q is-active getgo.service)
    then
    systemctl stop getgo
fi
mkdir -p /apps/getgo/bin/
cp -u ./bin/getgo /apps/getgo/bin/getgo
cp getgo.service /etc/systemd/system/getgo.service
systemctl daemon-reload
systemctl start getgo
echo "GetGo installed on your machine."
