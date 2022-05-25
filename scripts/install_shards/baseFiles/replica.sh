#!/bin/bash
set -e

mongo --port 10119 <<EOF
use admin
cfg={_id:"cfgReplica",members:[{_id:1,host:'127.0.0.1:10119'}]}
rs.initiate(cfg)
exit
EOF

echo waiting election...
sleep 3

mongo --port 10119 <<EOF
use admin
db.createUser({user:"root",pwd:"root2021_tsmp_isv_cn",roles:["root"]})
db.auth("root","root2021_tsmp_isv_cn")
db.createUser({user:"admin",pwd:"admin2021_tsmp_isv_cn",roles:[{role:"userAdminAnyDatabase",db:"admin"}]})
db.createUser({user:"tsmpadmin",pwd:"tsmpadmin2021_isv_cn",roles:[{role:"readWriteAnyDatabase",db:"admin"}]})
exit
EOF

mongo --port 10219 <<EOF
use admin
cfg={_id:"shardReplica",members:[{_id:1,host:'127.0.0.1:10219'}]}
rs.initiate(cfg)
EOF

echo waiting election...
sleep 3

mongo --port 10219 <<EOF
use admin
db.createUser({user:"root",pwd:"root2021_tsmp_isv_cn",roles:["root"]})
db.auth("root","root2021_tsmp_isv_cn")
db.createUser({user:"admin",pwd:"admin2021_tsmp_isv_cn",roles:[{role:"userAdminAnyDatabase",db:"admin"}]})
db.createUser({user:"tsmpadmin",pwd:"tsmpadmin2021_isv_cn",roles:[{role:"readWriteAnyDatabase",db:"admin"}]})
exit
EOF
