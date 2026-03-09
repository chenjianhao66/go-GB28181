# 用 Go 语言实现GB28181-2016标准的网络视频平台
go-GB28181是一个基于GB28181-2016标准实现的网络视频平台，用 Go 语言实现，实现了 SIP 协议和信令服务器。

# 功能特性

- [x] 信令服务：基于 GB28181-2016 标准实现完整的 SIP 信令处理
- [x] 设备管理：支持设备注册、注销、心跳检测、在线状态管理
- [x] 实时点播：支持向设备发起视频点播请求，获取实时视频流
- [x] 通道管理：支持设备通道查询、目录订阅与通知
- [x] 设备控制：支持云台 PTZ 控制、预置位操作
- [x] Web 管理界面：提供可视化界面进行设备管理和点播操作
- [x] 零依赖启动


# 快速开始

## Linux 环境

```bash
# 克隆源码
git clone https://github.com/chenjianhao66/go-GB28181.git
cd go-GB28181

# 编译（前端 + 后端）
make build

# 运行
./release/linux/go-GB28181 -c config/gbserver.yml
```

## Windows 环境

推荐直接下载 Release 压缩包，解压后运行可执行文件。



# 配置说明
| 配置项 | 是否必选 | 默认值 | 说明 |
|--------|----------|--------|------|
| **server** | | | |
| server.port | 是 | 18080 | HTTP 服务监听端口 |
| **sqlite** | | | |
| sqlite.path | 是 | ./db | 数据库文件存放目录 |
| sqlite.file | 是 | gbserver.db | SQLite 数据库文件名 |
| sqlite.username | 是 | root | 数据库用户名 |
| sqlite.password | 是 | root | 数据库密码 |
| sqlite.database | 是 | go-gb28181 | 数据库名称 |
| **nutsdb** | | | |
| nutsdb.path | 是 | ./db | NutsDB 数据存放目录 |
| **sip** | | | |
| sip.ip | 是 | 0.0.0.0 | 本机 IP 地址，用于接收设备连接 |
| sip.port | 否 | 5060 | SIP 服务监听端口 |
| sip.domain | 否 | 4401020049 | SIP 域，采用 ID 统一编码前十位 |
| sip.id | 否 | 44010200492000000001 | SIP 服务 ID |
| sip.password | 否 | - | 设备默认认证密码，移除则不校验 |
| sip.user-agent | 否 | gb | User-Agent 头信息 |
| **media** | | | |
| media.id | 是 | - | ZLMediaKit 服务器唯一 ID |
| media.ip | 是 | - | ZLMediaKit 服务器内网 IP |
| media.http-port | 是 | 80 | ZLMediaKit HTTP 端口 |
| media.secret | 否 | - | ZLMediaKit API 密钥 |
| **log** | | | |
| log.level | 否 | debug | 日志级别 (debug/info/warn/error) |
| log.path | 否 | ./log | 日志文件存放目录 |
| log.file | 否 | gbserver.log | 日志文件名 |
| log.maxSize | 否 | 1 | 单个日志文件最大大小 (MB) |
| log.maxBackups | 否 | 30 | 保留的旧日志文件最大数量 |
| log.maxAge | 否 | 30 | 日志文件保留最大天数 |





# 依赖服务

- **ZLMediaKit**：流媒体服务器，负责 RTP 视频流转发（必选）
- **SQLite**：持久化存储设备、通道等数据
- **NutsDB**：内存数据库，用于缓存 SIP 会话等热数据

# 项目目录结构

```
.
├── api                          ## 自动生成的 Swagger API 文档
│   └── swagger
├── cmd                          ## 程序入口
│   ├── gbctl                    ## 命令行客户端入口
│   └── gbserver                 ## GB28181 服务端入口
├── config                       ## 配置文件目录
├── db                           ## 数据库文件（SQLite、NutsDB）
├── docs                         ## 开发文档
│   ├── develop                  ## 开发指南
│   └── images                   ## 图片资源
├── internal                     ## 源代码
│   ├── config                   ## 配置模块
│   ├── gbctl                    ## 命令行客户端实现
│   ├── gbserver                 ## GB28181 服务端实现
│   │   ├── controller           ## HTTP 控制器
│   │   ├── gb                   ## GB28181 核心业务逻辑
│   │   ├── service              ## 业务服务层
│   │   ├── storage               ## 数据存储层
│   │   └── ui                   ## Web 前端页面
│   └── pkg                      ## 公共工具包
│       ├── gbsip                 ## SIP 协议实现
│       ├── model                ## 数据模型
│       └── parser               ## XML 解析
├── log                          ## 日志文件目录
├── release                      ## 编译产物
├── Makefile                     ## 编译脚本
└── README.md                    ## 项目说明
```

# 参考文档

对项目中有歧义的地方做了文档说明，请参考 `/docs` 目录


# 参考项目

流媒体服务基于@夏楚 [ZLMediaKit](https://github.com/ZLMediaKit/ZLMediaKit) 

国标处理逻辑基于[wvp-GB28181-pro](https://github.com/648540858/wvp-GB28181-pro) 

sip协议处理基于[go-sip](https://github.com/ghettovoice/gosip)
