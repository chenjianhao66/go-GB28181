# 用 Go 语言实现GB28181-2016标准的网络视频平台
go-GB28181是一个基于GB28181-2016标准实现的网络视频平台，用 Go 语言实现，实现了 SIP 协议和信令服务器。

# 开发计划

该项目还在积极开发中，下列是已经实现的和后续待实现的功能列表：

- [x] 注册和注销
- [x] 实时视音频点播
- [ ] 控制
  - [ ] 设备控制
    - [x] 云台控制
    - [ ] 录像控制
  - [x] 设备配置
- [ ] 信息查询
  - [x] 设备目录查询
  - [ ] 设备状态查询
  - [ ] 文件目录查询
  - [ ] 报警查询
  - [x] 设备配置查询
  - [x] 设备信息查询
- [ ] 通知 
  - [x] 状态信息报送（心跳）
  - [ ] 报警通知
  - [ ] 媒体通知
  - [ ] 语音广播通知
- [ ] 语音广播和语音对讲


# 项目目录结构

```
├── api                             ## 自动生成的接口文档
│   └── swagger                     
├── cmd                             ## 组件的main函数
│   ├── gbctl                       ## gb视频平台服务的启动函数
│   └── gbserver                    ## gb命令行客户端的启动函数
├── config                          ## 存放各个组件的配置文件，以组件名为文件名
│   ├── application-dev.yml
│   ├── gbctl.yml
│   └── gbserver.yml
├── docs                            ## 开发文档和用户文档
│   ├── develop
│   ├── guide
│   └── images
├── go.mod
├── go.sum
├── internal        
│   ├── config                      ## 废弃，后续删除
│   ├── gbctl                       ## gb命令行客户端的实现
│   ├── gbserver                    ## gb视频平台的实现
│   └── pkg                         ## 公共包
├── main.go
├── Makefile
└── README.md
```

# 参考文档

对项目中有歧义的地方做了文档说明，请参考 `/docs` 目录


# 参考项目

流媒体服务基于@夏楚 [ZLMediaKit](https://github.com/ZLMediaKit/ZLMediaKit) 

国标处理逻辑基于[wvp-GB28181-pro](https://github.com/648540858/wvp-GB28181-pro) 

sip协议处理基于[go-sip](https://github.com/ghettovoice/gosip)