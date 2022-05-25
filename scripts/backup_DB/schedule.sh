#!/bin/bash
curDir=$(cd "$(dirname $0)";pwd)
targetpath=$1
backupSaveMonths=$2
if [ ${#targetpath} -eq 0 ]
then
	echo "no path parameter specified, please specify target path!!!"
	exit 1
fi
if [[ $backupSaveMonths -eq 0 ]]
then
	echo "no save months parameter specified, please specify back up save months!!!"
	exit 1
fi
targetpath=${targetpath%*/}

allMongoCMD="22a\0 1 * * * root bash $curDir/mongodump_all.sh $targetpath $backupSaveMonths >/dev/null 2>&1"	#mongodb全量备份
allmysqlCMD="22a\0 1 * * * root bash $curDir/mysqldump_all.sh $targetpath $backupSaveMonths >/dev/null 2>&1"	#mysql全量备份
incMongoCMD="22a\0 4 * * * root bash $curDir/mongodump_inc.sh $targetpath $backupSaveMonths >/dev/null 2>&1"	#mongodb增量备份
incMysqlCMD="22a\0 4 * * * root bash $curDir/mysqldump_inc.sh $targetpath $backupSaveMonths >/dev/null 2>&1"	#mysql增量备份

sed -i "$allMongoCMD" /etc/crontab
sed -i "$allmysqlCMD" /etc/crontab
sed -i "$incMongoCMD" /etc/crontab
sed -i "$incMysqlCMD" /etc/crontab

if [ $? -eq 0 ]
then
	echo "cron successfully, delete myself, bye~"
	rm $0
fi
