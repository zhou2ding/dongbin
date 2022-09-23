#!/bin/bash
startTime=$(date)
echo "开始恢复：$startTime"

dbname='demo'
user='demo'
password='demo'
authdb='admin'

sourcefile=$1
if [ ${#sourcefile} -eq 0 ]
then
    echo "no source file specified, please specify source file!!!"
    exit 1
fi
sourcefile=${sourcefile%*/}

if [[ $sourcefile =~ ".sql" ]];then
	mysql -h$host -P$port -u$user -p$password $dbname < $sourcefile
elif [[ $sourcefile =~ "mysql-bin" ]];then
    mysqlbinlog --no-defaults $sourcefile | mysql -u$user -p$password
fi

endTime=$(date)
echo "恢复完成：$endTime"