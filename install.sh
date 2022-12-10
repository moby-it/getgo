#!/bin/bash

systemctl stop getgo
mkdir -p /apps/getgo/bin/
cp -u ./bin/getgo /apps/getgo/bin/getgo
cp getgo.service /etc/systemd/system/getgo.service
systemctl daemon-reload
systemctl start getgo