package main

import (
	_ "github.com/chenjianhao66/go-GB28181/api/swagger"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/server"
	"os"
	"os/signal"
	"syscall"
)

// @title           Go-GB28181项目前端APi接口
// @version         1.0
// @description     Go-GB28181是一个基于GB28181-2016标准实现的网络视频平台，用 Go 语言实现，实现了 SIP 协议和信令服务器。
// @termsOfService  http://swagger.io/terms/

// @contact.name   github homepage
// @contact.url    https://github.com/chenjianhao66/go-GB28181
// @contact.email  jianhao_c@qq.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:18080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	println(banner)
	ListenAndServeWithSignal()
}

var banner = `
  _____              _____ ____ ___   ___  __  ___  __ 
 / ____|            / ____|  _ \__ \ / _ \/_ |/ _ \/_ |
| |  __  ___ ______| |  __| |_) | ) | (_) || | (_) || |
| | |_ |/ _ \______| | |_ |  _ < / / > _ < | |> _ < | |
| |__| | (_) |     | |__| | |_) / /_| (_) || | (_) || |
 \_____|\___/       \_____|____/____|\___/ |_|\___/ |_|
`

func ListenAndServeWithSignal() {
	closeChan := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	// 监听并捕获 sighup（挂起）、sigquit(退出)、sigterm（终止）、sigint（终端）这些信号量，并将这些信号量写入到sigCh管道中
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	ListenAndServe(closeChan)
}

func ListenAndServe(closeChan <-chan struct{}) {
	s := server.NewServer()
	go func() {
		<-closeChan
		log.Info("gb server shutdown....")
		_ = s.Close()
	}()

	s.Run()
}
