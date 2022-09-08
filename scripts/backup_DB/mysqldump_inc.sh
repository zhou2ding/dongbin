#!/bin/bash
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

user='demo'
password='demo'
nowMonth=$(date "+%Y%m")
nowDay=$(date "+%Y%m%d")

if [ ! "$(ls $targetpath)" ];then
	#目标路径下为空，退出
	echo "Empty target directory, wait for all backup execute"
	exit 0
else
	#目标路径下非空，进一步检查mysql_all
	if ls $targetpath | grep "mysql_all" &>/dev/null;then
		#mysql_all存在，进一步检查mysql_inc
		if ls $targetpath | grep "mysql_inc" &>/dev/null;then
			#mysql_inc存在，非第一次增量备份
			echo "this is not first inc backup"
		else
			#mysql_inc不存在，是第一次增量备份
			echo "this is first inc backup"
		fi
	else
		#mysql_all不存在，退出
		echo "all backup not execute, wait for all backup execute"
		exit 0
	fi
fi

#判断当月的全量备份是否在当天完成的，是则退出脚本
allBackupDir=$targetpath/mysql_all/$nowMonth
if [[ -d $allBackupDir ]];then
	successFile=`ls $allBackupDir | grep "success_"`
	allBackupidx=`expr index $successFile "_"`
	allBackupDay=${successFile:allBackupidx}
	if [[ $allBackupDay = $nowDay ]];then
		echo "today has all backup, take s rest~~~"
		exit 0
	fi
fi

#判断当天是否有新增的binlog，没有则退出脚本
mysqladmin -u$user -p$password flush-logs #先刷新
binlogs=$(cat /var/log/mysql/mysql-bin.index)
arr=(${binlogs})
last=$((${#arr[@]}-1))
newestBin=${arr[$last]}
newestBin=${newestBin##*/}
if [[ -d $allBackupDir ]];then
	lastbinFile=`ls $allBackupDir | grep "lastbin_"`
	lastbinIdx=`expr index $lastbinFile "_"`
	allBackupLastBin=${lastbinFile:lastbinIdx}
	if [[ $allBackupLastBin = $newestBin ]];then
		echo "no new binlog, bye bye~~~"
		exit 0
	fi
fi

if [ ! -d $targetpath/mysql_inc/${nowDay} ]
then
	mkdir -p $targetpath/mysql_inc/${nowDay}
fi


#拷贝新出现的binlog
allBackupBinlogs=$(cat $targetpath/mysql_all/${nowMonth}/binlogs)
allBackupArr=(${allBackupBinlogs})
for binlog in ${arr[@]};do
	tempStr=$(echo "${allBackupArr[@]}" | grep -w ${binlog##*/})
	if [[ -z $tempStr ]];then
		echo "there is a new binlog:$binlog"
		cp $binlog $targetpath/mysql_inc/${nowDay}
	fi
done
