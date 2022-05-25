#!/bin/bash
set -e
# mongo --port 10001<<EOF
#use admin# db.createUser({user:"root",pwd:"root2021_tsmp_isv_com",roles:["root"]})
#db.auth("root","root2021_tsmp_isv_com")
# db.createUser({user:"admin",pwd:"admin2021_tsmp_isv_com",roles:[{role:"userAdminAnyDatabase",db:"admin"},{role:"clusterAdmin",db:"admin"}]})
#db.createUser({user:"tsmpadmin",pwd:"tsmpadmin2021_isv_com",roles:[{role:"readWriteAnyDatabase",db:"admin"}]})
# exit
#EOF

mongo --port 10001 <<EOF
use admin
db.auth("admin","admin2021_tsmp_isv_cn")
db.grantRolesToUser("admin",[{role:"clusterAdmin",db:"admin"}])
sh.addShard("shardReplica/127.0.0.1:10219")
sh.enableSharding("tsmp")
sh.shardCollection("tsmp.train_records",{pd:1,it:1})
sh.shardCollection("tsmp.detections_360",{pd:1,it:1})
sh.shardCollection("tsmp.faults_360",{pd:1,it:1})
sh.shardCollection("tsmp.measurements_twd",{pd:1,it:1})
sh.shardCollection("tsmp.faults_twd",{pd:1,it:1})
sh.shardCollection("tsmp.faults_twd2d",{pd:1,it:1})
sh.shardCollection("tsmp.detections_panto",{pd:1,it:1})
sh.shardCollection("tsmp.faults_panto",{pd:1,it:1})
exit
EOF
