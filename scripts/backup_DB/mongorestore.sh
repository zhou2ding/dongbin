#!/bin/bash
startTime=$((`date -d "$nowStr" "+%s"`*1000+`date -d "$nowStr" "+%N"`/1000000))

host='127.0.0.1'
port='10005'
dbname='zdb'
user='zdbadmin'
password='zdbadmin2021_2ding_cn'
authdb='admin'

restoreFunc() {
	targetpath=$1
	if [ ${#targetpath} -eq 0 ]
	then
		echo "no path parameter specified, please specify target path!!!"
		exit 1
	fi
	targetpath=${targetpath%*/}
	 
	mongorestore -h ${host}:${port} -d $dbname -u $user -p $password --authenticationDatabase $authdb ${targetpath}
}

restoreFunc $1

endTime=$((`date -d "$nowStr" "+%s"`*1000+`date -d "$nowStr" "+%N"`/1000000))
costTime=`expr $endTime - $startTime`
curDir=$(cd "$(dirname $0)";pwd)
touch $curDir/systemTest_mongoRestore_costs_$costTime
