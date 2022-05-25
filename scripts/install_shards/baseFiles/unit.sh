#!/bin/bash
dir=$(cd "$(dirname $0)";pwd)
for d in `ls`
do
	if [ -d $d ]
	then
		if [[ $d =~ "shard" || $d =~ "configsvr" || $d =~ "router" ]]
		then
			bash $dir/$d/del.sh
		fi
	fi
done

rm -rf /var/lib/mongodb-shards