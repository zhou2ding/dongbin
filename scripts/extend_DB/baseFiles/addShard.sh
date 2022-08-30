#!/bin/bash
set -e

mongo --host SPECIFIED_IP --port SPECIFIED_PORT <<EOF
use admin
cfg={_id:"SHARD_NAME",members:[{_id:1,host:'SPECIFIED_IP:10220'}]}
rs.initiate(cfg)
exit
EOF

echo waiting election...
sleep 5

mongo --host SPECIFIED_IP --port SPECIFIED_PORT <<EOF
use admin
db.createUser({user:"root",pwd:"root2021_zdb_2ding_cn",roles:["root"]})
db.auth("root","root2021_zdb_2ding_cn")
db.createUser({user:"admin",pwd:"admin2021_zdb_2ding_cn",roles:[{role:"userAdminAnyDatabase",db:"admin"}]})
db.createUser({user:"zdbadmin",pwd:"zdbadmin2021_2ding_cn",roles:[{role:"readWriteAnyDatabase",db:"admin"}]})
exit
EOF

mongo --port 10001 <<EOF
use admin
db.auth("admin","admin2021_zdb_2ding_cn")
sh.addShard("SHARD_NAME/SPECIFIED_IP:SPECIFIED_PORT")
exit
EOF