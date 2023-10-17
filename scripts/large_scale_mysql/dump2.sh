#!/bin/bash

# 设置变量
output_dir="/home/zdb/sqldumps"
batch_size=1000000  # 指定每批次的大小

# 表名数组
TABLE_NAMES=(
    "demo_table"
)

# 创建输出目录
mkdir -p "$output_dir"

# 遍历所有表名
for table_name in "${TABLE_NAMES[@]}"; do
    # 获取表的总行数
    total_rows=$(/opt/MAIOT/MAIOT/deps/mysql/bin/mysql -h127.0.0.1 -u root -p'my_pass' -N -s -e "SELECT TABLE_ROWS FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='data' AND TABLE_NAME='$table_name';")

    # 计算批次数
    num_batches=$((total_rows / batch_size + 1))

    # 导出每个批次的数据
    for ((i = 0; i < num_batches; i++)); do
        offset=$((i * batch_size))
        output_csv="$output_dir/${table_name}_batch_${i}.csv"
        output_sql="$output_dir/${table_name}_batch_${i}.sql"

        # 导出数据为 CSV 文件
        /opt/MAIOT/MAIOT/deps/mysql/bin/mysql -h127.0.0.1 -u root -p'my_pass' -N -s -e "SELECT * FROM data.$table_name LIMIT $offset, $batch_size" data > "$output_csv"

        # 转换 CSV 文件为 SQL 文件
        awk -F'\t' '{printf "INSERT INTO '"$table_name"' VALUES ("; for(i=1; i<=NF; ++i) printf i<NF ? "\""$i"\"," : "\""$i"\""; print ");"}' "$output_csv" > "$output_sql"

        echo "表 $table_name 的批次 $i 导出完成: $output_sql"
    done
done

echo "导出完成"
