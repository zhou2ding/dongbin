#!/bin/bash

dir=$(cd "$(dirname $0)";pwd)

# 初始化配置服务器和分片
echo waiting for service start...
for d in `ls`
do
if [ -d $d ]
then
if [[ $d =~ "shard" || $d =~ "configsvr" ]]
then
bash $dir/$d/mongod_service.sh
fi
fi
done

# 构建分片和配置服务器的副本集
sleep 3
bash $dir/scripts/replica.sh

# 初始化路由服务器
echo waiting for service start...
for d in `ls`
do
if [ -d $d ]
then
if [[ $d =~ "router" ]]
then
bash $dir/$d/mongos_service.sh
fi
fi
done

# 对数据库进行分片
sleep 3
echo waiting for service start...
bash $dir/scripts/shards.sh
#sleep 3
#systemctl status |grep mongo
