#!/bin/bash

curDir=$(cd $(dirname $0);pwd)
now=$(date +"%Y%m%d-%H")
mkdir -p $curDir/$now
selfpid=$$  #脚本自身pid

#方式1
top -p 1368,3429132,1376,1377,$selfpid -n 2100 -d 1 -b > ${curDir}/${now}/${serviceName}

#方式2
declare -A services #key是进程的pid，value是进程的名称
services["1368"]="mongod"
services["3429132"]="mysqld"
services["1376"]="rabbitmq-server"
services["1377"]="redis-server"

for pid in ${!services[*]};do
    {
        serviceName=${services[$pid]}
        top -p ${pid} -n 2100 -d 1 -b > ${curDir}/${now}/${serviceName}
    }&
done
wait
