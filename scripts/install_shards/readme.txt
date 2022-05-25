1. 首次安装，直接执行environment下的install.sh以初始化安装环境
2. 非首次安装，
需要 openssl ，如果没有，进入 0.9.0\environment\packages\1\openssl  执行  sudo dpkg -i openssl_1.1.1f-1ubuntu2.12_amd64.deb
已有mongod的单例服务，直接执行environment\mongodb-shards下的install.sh以把单例换成集群
3. 卸载集群，去/var/lib/mongodb-shards下执行unit.sh（会把数据和日志也清空）