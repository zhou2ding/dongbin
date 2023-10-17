#!/bin/bash

# 数据库连接信息
DB_USER="root"
DB_PASSWORD='my_pass'
DB_NAME="demo_db"

# 通过参数指定表名
TABLE_NAME=$1

if [ -z "$TABLE_NAME" ]; then
    echo "Usage: $0 <table_name>"
    exit 1
fi

# 遍历所有sql文件
for sql_file in *.sql; do
    # 导入数据到数据库
    /opt/MAIOT/MAIOT/deps/mysql/bin/mysql -h127.0.0.1 -u"$DB_USER" -p"$DB_PASSWORD" -D"$DB_NAME" < "$sql_file"
done
