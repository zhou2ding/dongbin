#!/bin/bash

curDir=$(cd "$(dirname $0)";pwd)
prefix=${curDir##*/}

service_script="
[Unit]
Description=An object/document-oriented database
Documentation=man:mongod(1)
After=network.target

[Service]
RuntimeDirectory=mongodb
RuntimeDirectoryMode=0755
EnvironmentFile=-/etc/default/mongodb
Environment=CONF=CUR_DIR/${prefix}.conf
Environment=SOCKETPATH=/run/mongodb
ExecStart=/usr/bin/mongos -f CUR_DIR/${prefix}.conf
LimitFSIZE=infinity
LimitCPU=infinity
LimitAS=infinity
LimitNOFILE=64000
LimitNPROC=64000

[Install]
WantedBy=multi-user.target
"
service_name=mongodb_${prefix}.service
echo "$service_script" | sed "s|CUR_DIR|${curDir}|g" > /usr/lib/systemd/system/$service_name

systemctl enable $service_name
systemctl start $service_name
