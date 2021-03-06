#!/bin/bash
host='127.0.0.1'
port='10005'
dbname='tsmp'
user='tsmpadmin'
password='tsmpadmin2021_isv_cn'
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