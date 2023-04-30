package gbserver

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/app"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
)

const description = `基于GB28181-2016标准实现的网络视频平台，用 Go 语言实现，实现了 SIP 协议和信令服务器。

后续还有其他功能等待支持，如果该程序有帮到你的话请去仓库给作者点一个Start吧~
	https://github.com/chenjianhao66/go-GB28181
`

const banner = `
  _____              _____ ____ ___   ___  __  ___  __ 
 / ____|            / ____|  _ \__ \ / _ \/_ |/ _ \/_ |
| |  __  ___ ______| |  __| |_) | ) | (_) || | (_) || |
| | |_ |/ _ \______| | |_ |  _ < / / > _ < | |> _ < | |
| |__| | (_) |     | |__| | |_) / /_| (_) || | (_) || |
 \_____|\___/       \_____|____/____|\___/ |_|\___/ |_|
`

func NewApp(basename string) *app.App {
	option := newGbOption()
	return app.NewApp("国标信令服务",
		basename,
		app.WithOption(option),
		app.WithBanner(banner),
		app.WithDescription(description),
		app.WithRunFunc(run(option)),
	)
}

func run(opt *GbOption) app.RunFunc {
	return func(basename string) error {
		log.Init(opt.LogOption)
		return NewServer(opt).Run()
	}
}
