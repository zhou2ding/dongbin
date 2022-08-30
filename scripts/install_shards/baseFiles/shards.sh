#!/bin/bash
set -e
# mongo --port 10001<<EOF
#use admin# db.createUser({user:"root",pwd:"root2021_zdb_2ding_com",roles:["root"]})
#db.auth("root","root2021_zdb_2ding_com")
# db.createUser({user:"admin",pwd:"admin2021_zdb_2ding_com",roles:[{role:"userAdminAnyDatabase",db:"admin"},{role:"clusterAdmin",db:"admin"}]})
#db.createUser({user:"zdbadmin",pwd:"zdbadmin2021_2ding_com",roles:[{role:"readWriteAnyDatabase",db:"admin"}]})
# exit
#EOF

mongo --port 10001 <<EOF
use admin
db.auth("admin","admin2021_zdb_2ding_cn")
db.grantRolesToUser("admin",[{role:"clusterAdmin",db:"admin"}])
sh.addShard("shardReplica/127.0.0.1:10219")
sh.enableSharding("zdb")
sh.shardCollection("zdb.train_records",{pd:1,it:1})
sh.shardCollection("zdb.detections_360",{pd:1,it:1})
sh.shardCollection("zdb.faults_360",{pd:1,it:1})
sh.shardCollection("zdb.measurements_twd",{pd:1,it:1})
sh.shardCollection("zdb.faults_twd",{pd:1,it:1})
sh.shardCollection("zdb.faults_twd2d",{pd:1,it:1})
sh.shardCollection("zdb.detections_panto",{pd:1,it:1})
sh.shardCollection("zdb.faults_panto",{pd:1,it:1})
exit
EOF
