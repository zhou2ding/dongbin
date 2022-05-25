#!/bin/bash
curDir=$(cd "$(dirname $0)";pwd)
prefix=${curDir##*/}

service_name=mongodb_${prefix}.service
echo "start to uninstall $service_name"
systemctl stop $service_name
systemctl disable $service_name
rm /usr/lib/systemd/system/$service_name

