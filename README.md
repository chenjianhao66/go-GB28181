# 用 Go 语言实现GB28181-2016标准的网络视频平台
go-GB28181是一个基于GB28181-2016标准实现的网络视频平台，用 Go 语言实现，实现了 SIP 协议和信令服务器。

# 开发计划

该项目还在积极开发中，下列是已经实现的和后续待实现的功能列表：

- [x] 注册和注销
- [x] 实时视音频点播
- [ ] 控制
  - [ ] 设备控制
  - [ ] 设备配置
- [ ] 信息查询
  - [ ] 设备目录查询
  - [ ] 设备状态查询
  - [ ] 文件目录查询
  - [ ] 报警查询
  - [ ] 设备配置查询
  - [x] 设备信息查询
- [ ] 通知 
  - [x] 状态信息报送（心跳）
  - [ ] 报警通知
  - [ ] 媒体通知
  - [ ] 语音广播通知
- [ ] 语音广播和语音对讲


# 项目目录结构

```
├── config                    ## 配置文件目录以及模板
│   ├── application-dev.yml
│   └── application.yml
├── go.mod
├── go.sum
├── internal                  ## 私有目录
│   ├── config                ## 配置文件
│   ├── controller            ## api控制器
│   ├── gb                    ## 国标
│   ├── log                   ## log包
│   ├── model                 ## 数据库实体struct
│   ├── parser                ## xml、json等解析代码包
│   ├── server                ## 集成api、sip等服务
│   ├── service               ## mvc中的service
│   └── store                 ## mvc中的dao层
├── main.go
├── Makefile
└── README.md
```


# 参考项目

流媒体服务基于@夏楚 [ZLMediaKit](https://github.com/ZLMediaKit/ZLMediaKit) 

国标处理逻辑基于[wvp-GB28181-pro](https://github.com/648540858/wvp-GB28181-pro) 

sip协议处理基于[go-sip](https://github.com/ghettovoice/gosip)