package gbctl

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/app"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

const description = `这是一个实现了国标标准的模拟摄像头，它将实现国标的功能用于调试。
下面是已经实现了的功能：
	1. 注册与注销功能

后续还有其他功能等待支持，如果该程序有帮到你的话请去仓库给作者点一个Start吧~
	https://github.com/chenjianhao66/go-GB28181
`

func NewApp(basename string) *app.App {
	option := newCTLOption()
	return app.NewApp("国标命令行客户端",
		basename,
		app.WithOption(option),
		app.WithDescription(description),
		app.WithRunFunc(run(option)),
	)
}

func run(opt *ctlOption) app.RunFunc {
	return func(basename string) error {
		log.Init(opt.LogOption)
		log.Info("exec gbctl success....")
		log.Info(cast.ToString(viper.GetString("sip.id")))
		viper.New()
		strings := viper.GetStringSlice("client.channel.list")

		log.Info(strings)
		return nil
	}
}
