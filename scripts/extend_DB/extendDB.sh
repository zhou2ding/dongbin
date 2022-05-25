#!/bin/bash
set -e
curDir=$(cd "$(dirname $0)";pwd)

ip=$1
if [ ${#ip} -eq 0 ]
then
	echo IP is not specified, please specify IP !!!
	exit 1
fi
chmod -R +x $curDir

(
	cd $curDir/baseFiles
	for d in `ls`
	do
		if [ -d $d ]
		then
			if [[ $d =~ "shard" ]]
			then
				# 把新分片的配置文件和启动脚本拷贝到集群文件夹下面
				mkdir -p /var/lib/mongodb-shards/$d/logs
				mkdir -p /var/lib/mongodb-shards/$d/data
				cp ./$d/*.conf /var/lib/mongodb-shards/$d
				cp ./$d/mongod_service.sh /var/lib/mongodb-shards/$d
				cp /var/lib/mongodb-shards/router1/mongo.keyfile /var/lib/mongodb-shards/$d

				# 启动新分片的mongod服务
				echo waiting mongodb service start...
				bash /var/lib/mongodb-shards/$d/mongod_service.sh
				sleep 3

				# 根据指定的IP、端口和从配置文件中读到的replSetName的值初始化新分片
				port=`awk 'NR==14{print $2}' ./$d/*.conf`
				name=`awk 'NR==9{print $2}' ./$d/*.conf`
				cp ${curDir}/baseFiles/addShard.sh ${curDir}/baseFiles/addShard_temp.sh
				sed -i "s|SHARD_NAME|$name|g" ${curDir}/baseFiles/addShard_temp.sh
				sed -i "s|SPECIFIED_IP|$ip|g" ${curDir}/baseFiles/addShard_temp.sh
				sed -i "s|SPECIFIED_PORT|$port|g" ${curDir}/baseFiles/addShard_temp.sh
				bash ${curDir}/baseFiles/addShard_temp.sh
				rm ${curDir}/baseFiles/addShard_temp.sh
			fi
		fi
	done
)
