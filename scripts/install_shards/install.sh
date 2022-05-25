#!/bin/bash
# 先停掉已有的mongod.service
systemctl stop mongod.service
systemctl disable mongod.service
sleep 2

# 创建文件夹、拷贝配置文件、拷贝配置副本集和分片的脚本
curDir=$(cd "$(dirname $0)";pwd)
chmod +x $curDir/baseFiles/*.sh
mkdir /var/lib/mongodb-shards
mkdir /var/lib/mongodb-shards/configsvr1
mkdir /var/lib/mongodb-shards/shard1
mkdir /var/lib/mongodb-shards/router1
mkdir /var/lib/mongodb-shards/configsvr1/logs
mkdir /var/lib/mongodb-shards/shard1/logs
mkdir /var/lib/mongodb-shards/router1/logs
mkdir /var/lib/mongodb-shards/configsvr1/data
mkdir /var/lib/mongodb-shards/shard1/data
mkdir /var/lib/mongodb-shards/scripts
cp $curDir/baseFiles/configsvr1.conf /var/lib/mongodb-shards/configsvr1
cp $curDir/baseFiles/shard1.conf /var/lib/mongodb-shards/shard1
cp $curDir/baseFiles/router1.conf /var/lib/mongodb-shards/router1

cp $curDir/baseFiles/init.sh /var/lib/mongodb-shards/
cp $curDir/baseFiles/unit.sh /var/lib/mongodb-shards/
cp $curDir/baseFiles/del.sh /var/lib/mongodb-shards/configsvr1
cp $curDir/baseFiles/del.sh /var/lib/mongodb-shards/shard1
cp $curDir/baseFiles/del.sh /var/lib/mongodb-shards/router1
cp $curDir/baseFiles/mongod_service.sh /var/lib/mongodb-shards/configsvr1
cp $curDir/baseFiles/mongod_service.sh /var/lib/mongodb-shards/shard1
cp $curDir/baseFiles/mongos_service.sh /var/lib/mongodb-shards/router1
cp $curDir/baseFiles/replica.sh /var/lib/mongodb-shards/scripts
cp $curDir/baseFiles/shards.sh /var/lib/mongodb-shards/scripts

# 生成并拷贝秘钥文件
openssl rand -out $curDir/baseFiles/mongo.keyfile -hex 128
chmod 400 $curDir/baseFiles/mongo.keyfile
cp $curDir/baseFiles/mongo.keyfile /var/lib/mongodb-shards/configsvr1
cp $curDir/baseFiles/mongo.keyfile /var/lib/mongodb-shards/shard1
cp $curDir/baseFiles/mongo.keyfile /var/lib/mongodb-shards/router1

# 初始化分片集群
(
cd /var/lib/mongodb-shards
bash ./init.sh
)
