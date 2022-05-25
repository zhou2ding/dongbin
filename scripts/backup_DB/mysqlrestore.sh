#!/bin/bash
host='127.0.0.1'
port='10008'
dbname='tsmp'
user='root'
password='root2021@tsmp_isv_cn'
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
