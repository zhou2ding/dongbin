# dongbin
## 项目结构说明

```
|── bfhummingbird：前端工程
|── cmd
|     |--- user_manager: 用户管理服务
|     |--- blog_manager: 帖子管理服务
|── config：配置文件
|—— controller: 业务处理层
|── model: 对象模型      
|—— handler: 接口访问层
|── pkg: 公共模块
|     |── auth：鉴权加密模块
|     |── cache：缓存模块
|     |── mq：消息队列
|     |── nosql: 非关系型数据库模块
|     |── filemanager：文件读写模块
|     |── sql：关系型数据库模块
|     |── rpcserer: rpc通信模块
|     |── rpcservice: rpc服务模块
|     |── rpcclient: rpc客户端
|     |── license: 授权认证模块
|     |── util: 基础模块，包括日志、viper
|     |── wsService: websocket服务模块
|—— rpchandler： RPC接口访问层
|—— wshandler： websocket接口访问层
|── tool: 工具类
|── scripts: shell脚本
```
