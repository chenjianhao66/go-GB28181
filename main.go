package main

import (
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ListenAndServeWithSignal()
}

func ListenAndServeWithSignal() {
	// 声明关闭通知通道、系统信号量的通道
	closeChan := make(chan struct{})
	sigCh := make(chan os.Signal)
	// 监听并捕获 sighup（挂起）、sigquit(退出)、sigterm（终止）、sigint（终端）这些信号量，并将这些信号量写入到sigCh管道中
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	// 起goroutine来监听 sigCh管道的数据，有数据的就代表要退出程序了，往 closeChan 管道内发送数据
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
	s.Run()

	<-closeChan
	log.Info("close....")
}
