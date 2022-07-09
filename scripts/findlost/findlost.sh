#!/bin/bash

#用tar分包压缩后，文件名称的规律是：filename.tar.gzaa~filename.tar.gzyz，yz之后不是za，而是zaaa~zyzz
package=()
idx=0
for ((i=97;i<=121;i++));do
	t1=`echo $i | awk '{printf("%c", $1)}'`
	for ((j=97;j<=122;j++));do
		t2=`echo $j | awk '{printf("%c", $1)}'`
		package[$idx]="filename.tar.gz"$t1$t2
		idx=`expr $idx + 1`
	done
done

existPackage=()
idx2=0
for file in `ls`;do
	existPackage[$idx2]=$file
	idx2=`expr $idx2 + 1`
done

for elem in ${package[@]};do
	if echo "${existPackage[@]}" | grep -w ${elem##*/}  &>/dev/null;then
		x=0
	else
		echo missing:$elem
	fi
done
