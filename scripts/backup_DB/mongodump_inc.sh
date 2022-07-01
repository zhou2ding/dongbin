#!/bin/bash
nowStr=`date "+%Y-%m-%d %H:%M:%S"`     #当前时间
endTime=$((`date -d "$nowStr" "+%s"`*1000+`date -d "$nowStr" "+%N"`/1000000)) #毫秒级时间戳
beginTime=0

targetpath=$1
backupSaveMonths=$2
if [ ${#targetpath} -eq 0 ];then
	echo "no path parameter specified, please specify target path!!!"
	exit 1
fi
if [[ $backupSaveMonths -eq 0 ]];then
	echo "no save months parameter specified, please specify back up save months!!!"
	exit 1
fi
targetpath=${targetpath%*/}

host='127.0.0.1'
port='10005'
dbname='tsmp'
user='tsmpadmin'
password='tsmpadmin2021_isv_cn'
authdb='admin'
nowDay=$(date "+%Y%m%d")
nowMonth=$(date "+%Y%m")

function getTimeStap() {
	local p=$1
	if [ $p = "mongo_all" ];then
		d1=$(date "+%Y%m")
		for f in `ls $targetpath/$p/$d1/tsmp`;do
			if [[ $f =~ "timestap" ]];then
				ret=$f
				echo $ret
				return $?
			fi
		done
	elif [ $p = "mongo_inc" ];then
		d1=$2
		for f in `ls $targetpath/$p/$d1/tsmp`;do
			if [[ $f =~ "timestap" ]];then
				ret=$f
				echo $ret
				return $?
			fi
		done
	fi
}

isFirst="true"
if [ ! "$(ls $targetpath)" ];then
	#目标路径下为空，退出
	echo "Empty target directory, wait for all backup execute"
	exit 0
else
	#目标路径下非空，进一步检查mongo_all
	if ls $targetpath | grep "mongo_all" &>/dev/null;then
		#mongo_all存在，进一步检查mongo_inc
		if ls $targetpath | grep "mongo_inc" &>/dev/null;then
			#mongo_inc存在，非第一次增量备份，起点时间戳是最近一次增量备份的时间戳
			isFirst="false"
		else
			#mongo_inc不存在，是第一次增量备份，起点时间戳是第一次全量备份的时间戳
			timestap_all=$(getTimeStap mongo_all)
			idx1=`expr index $timestap_all "_"`
			beginTime=${timestap_all:idx1}
			echo "this is first inc backup, beginTime:"$beginTime
		fi
	else
		#mongo_all不存在，退出
		echo "all backup not execute, wait for all backup execute"
		exit 0
	fi
fi

recentTimeStap=0
i=0
if [[ $isFirst = "false" ]];then
	#查找最近一次的增量备份
	while [[ $recentTimeStap -eq 0 ]];do
		day=$((`date -d "-${i} day" "+%Y%m%d"`))
		timestap_inc=$(getTimeStap mongo_inc $day)
		idx2=`expr index $timestap_inc "_"`
		recentTimeStap=${timestap_inc:idx2}
		beginTime=$recentTimeStap
		let i++
		if [[ $i -gt 100 ]];then
			echo "inc backup is too little!!!"
			break
		fi
	done
fi

#判断当月的全量备份中，日期文件和当天是否一致，一致则退出
allBackupDir=$targetpath/mongo_all/$nowMonth/tsmp
if [[ -d $allBackupDir ]];then
	successFile=`ls $allBackupDir | grep "success_"`
	allBackupidx=`expr index $successFile "_"`
	allBackupDay=${successFile:allBackupidx}
	if [[ $allBackupDay = $nowDay ]];then
		echo "today has all backup, take s rest~~~"
		exit 0
	fi
fi

echo "filter:" '{"it":{"$gte":'$beginTime',"$lte":'$endTime'}}'
if [ $beginTime -ne 0 ];then
	if [[ ! -d ${targetpath}/mongo_inc/${nowDay} ]]
	then
		mkdir -p ${targetpath}/mongo_inc/${nowDay}
	fi
	collections=("train_records" "detections_panto" "faults_panto" "detections_360" "faults_360" "measurements_twd" "faults_twd" "faults_twd2d" "simple_train_records" "sim_train_records")
	for collection in ${collections[@]}
	do
		mongodump -h ${host}:${port} -d $dbname -o ${targetpath}/mongo_inc/${nowDay} -u $user -p $password --authenticationDatabase $authdb -c $collection -q '{"it":{"$gte":'$beginTime',"$lte":'$endTime'}}'
		if [ $? -eq 0 ]
		then
			nowTime=$((`date -d "$nowStr" "+%s"`*1000+`date -d "$nowStr" "+%N"`/1000000))	#毫秒级时间戳
			touch $targetpath/mongo_inc/${nowDay}/tsmp/timestap_${nowTime}
		else
			rm -rf $targetpath/mongo_inc/${nowDay}
		fi
	done
fi

#清理过期备份
backupFileCnt=`ls $targetpath/mongo_all -l |grep "^d"|wc -l`
if [[ $backupFileCnt -gt $backupSaveMonths ]];then
	monArr=()
	i=0
	for mon in `ls $targetpath/mongo_all`;do
		monArr[$i]=$mon
		i=`expr $i + 1`
	done

	for ((j=0;j<${#monArr[@]};j++));do
		diff=`expr ${#monArr[@]} - $backupSaveMonths`
		if [[ $j -lt $diff ]];then
			echo "clear expire backup files: $targetpath/mongo_all/${monArr[$j]}"
			rm -rf $targetpath/mongo_all/${monArr[$j]}
		fi
	done
fi
