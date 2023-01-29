package main

import (
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/server"
	"os"
	"os/signal"
	"syscall"
)

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
	sigCh := make(chan os.Signal)
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
