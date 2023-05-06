package gbserver

import (
	"context"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/gb"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage/cache"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	sip       *gb.Server
	apiServer *apiServer
	ctx       context.Context
	cancel    context.CancelFunc
	opt       *GbOption
}

func NewServer(opt *GbOption) *Server {
	ctx, cancelFunc := context.WithCancel(context.Background())
	apiConfig := &apiConfig{
		mediaOption:  opt.MediaOption,
		serverOption: opt.ServerOption,
		mysqlOption:  opt.MysqlOption,
	}
	gbConfig := &gb.SipConfig{
		SipOption:   opt.Sip,
		MysqlOption: opt.MysqlOption,
	}
	return &Server{
		sip:       gb.NewServer(gbConfig),
		apiServer: newApiServer(apiConfig),
		ctx:       ctx,
		cancel:    cancelFunc,
		opt:       opt,
	}
}

func (s *Server) Run() error {
	closeCh := listenAndServeWithSignal()
	go func() {
		<-closeCh
		if err := s.Close(); err != nil {
			log.Error(err)
		}
		log.Info("gbserver shutdown....")
	}()
	cache.InitCache(s.opt.RedisOption)
	s.apiServer.initRoute()
	eg, ctx := errgroup.WithContext(s.ctx)
	defer ctx.Done()

	eg.Go(func() error {
		log.Infof("bind: %s,start listening...", s.apiServer.h.Addr)
		if err := s.apiServer.h.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		return nil
	})

	eg.Go(func() error {
		return s.sip.ListenTCP()
	})

	eg.Go(func() error {
		return s.sip.ListenUDP()
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	if err := s.sip.Close(); err != nil {
		return err
	}
	cancel := s.cancel
	cancel()
	if err := s.apiServer.Close(); err != nil {
		return err
	}
	return nil
}

func listenAndServeWithSignal() <-chan struct{} {
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
	return closeChan
}
