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

#先清理过期备份
backupFileCnt=`ls $targetpath/mysql_all -l |grep "^d"|wc -l`
for ((i=$backupSaveMonths; i<=$backupFileCnt; i++));do
	delMonth=$((`date -d "-${i} month" "+%Y%m"`))
	echo "delete expired backup files: $targetpath/mysql_all/$delMonth"
	rm -rf $targetpath/mysql_all/$delMonth
done

host='127.0.0.1'
port='10008'
dbname='tsmp'
user='root'
password='root2021@tsmp_isv_cn'
nowMonth=$(date "+%Y%m")
nowDay=$(date "+%Y%m%d")

exist=0
function traverseDir() {
	for elem in `ls $1`
	do
		dof=$1"/"$elem
		if [[ -f $dof && $dof =~ "tsmp.sql" ]]
		then
			exist=`expr $exist + 1`
		fi
		if [[ -f $dof && $dof =~ "success" ]]
		then
			exist=`expr $exist + 1`
		fi
	done
}
#检查当月有没有备份，有则退出脚本
traverseDir $targetpath/mysql_all/${nowMonth}
if [ $exist == 2 ]
then
	echo "MySQL BackUp has done in this month"
	exit 0
fi

if [ ! -d $targetpath/mysql_all/${nowMonth} ]
then
	mkdir -p $targetpath/mysql_all/${nowMonth}
fi

#刷新binlog后执行备份
mysqladmin -h$host -P$port -u$user -p$password flush-logs
mysqldump -h$host -P$port -u$user -p$password $dbname > $targetpath/mysql_all/${nowMonth}/tsmp.sql

if [ $? -eq 0 ]
then
	echo "MySQL BackUp Successful"
	touch $targetpath/mysql_all/${nowMonth}/success_${nowDay}
	cp /var/log/mysql/mysql-bin.index $targetpath/mysql_all/${nowMonth}/binlogs
	chmod 644 $targetpath/mysql_all/${nowMonth}/binlogs
	binlogs=$(cat /var/log/mysql/mysql-bin.index)
	arr=(${binlogs})
	last=$((${#arr[@]}-1))
	lastbin=${arr[$last]}
	lastbin=${lastbin##*/}
	touch $targetpath/mysql_all/${nowMonth}/lastbin_${lastbin}
	nowStr=`date "+%Y-%m-%d %H:%M:%S"`
	nowTime=$((`date -d "$nowStr" "+%s"`*1000+`date -d "$nowStr" "+%N"`/1000000))	#毫秒级时间戳
	touch $targetpath/mysql_all/${nowMonth}/timestap_${nowTime}
else
	echo "MySQL BackUp Fail"
fi