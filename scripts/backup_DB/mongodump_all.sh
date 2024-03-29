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
	echo "no save days parameter specified, please specify back up save days!!!"
	exit 1
fi
targetpath=${targetpath%*/}
host='127.0.0.1'
port='27017'
dbname='demo'
user='demo'
password='demo'
authdb='admin'
nowMonth=$(date "+%Y%m")
nowDay=$(date "+%Y%m%d")

#当前是全量设备的数据，后续可配置设备
files=(
	"demo.bson"
	)

exist=0
function traverseDir() {
	for elem in `ls $1`;do
		dof=$1"/"$elem
		if [[ -f $dof && $dof =~ "success" ]];then
			exist=`expr $exist + 1`
		fi
		if [ -f $dof ];then
			if echo "${files[@]}" | grep -w ${dof##*/}  &>/dev/null;then
				exist=`expr $exist + 1`
			fi
		fi
	done
}
#检查当月有没有备份
traverseDir $targetpath/mongo_all/${nowMonth}/zdb
if [ $exist == 19 ]
then
	echo "MongoDB BackUp has done in this month"
	exit 0
fi

if [ ! -d ${targetpath}/mongo_all/${nowMonth} ]
then
	mkdir -p ${targetpath}/mongo_all/${nowMonth}
fi

#执行备份
mongodump -h ${host}:${port} -d $dbname -o ${targetpath}/mongo_all/${nowMonth} -u $user -p $password --authenticationDatabase $authdb
if [ $? -eq 0 ]
then
	if [ -d $targetpath/mongo_all/${nowMonth}/zdb ]
	then
		echo "MongoDB BackUp Successful"
		touch $targetpath/mongo_all/${nowMonth}/zdb/success_${nowDay}
		nowStr=`date "+%Y-%m-%d %H:%M:%S"`
		nowTime=$((`date -d "$nowStr" "+%s"`*1000+`date -d "$nowStr" "+%N"`/1000000))	#毫秒级时间戳
		touch $targetpath/mongo_all/${nowMonth}/zdb/timestap_${nowTime}
	else
		echo "MongoDB BackUp Fail"
	fi
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
			echo "clear expire backup files: $targetpath/mongo_inc/${monArr[$j]}"
			rm -rf $targetpath/mongo_all/${monArr[$j]}
			rm -rf $targetpath/mongo_inc/${monArr[$j]}*
		fi
	done
fi