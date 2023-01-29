package server

import (
	"context"
	"github.com/chenjianhao66/go-GB28181/internal/gb"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type Server struct {
	sip       *gb.Server
	apiServer *apiServer
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewServer() *Server {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Server{
		sip:       gb.NewServer(),
		apiServer: newApiServer(),
		ctx:       ctx,
		cancel:    cancelFunc,
	}
}

func (s *Server) Run() {
	s.apiServer.initRoute()
	eg, _ := errgroup.WithContext(s.ctx)

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
		panic(err)
	}

}

func (s *Server) Close() error {
	_ = s.sip.Close()
	_ = s.apiServer.Close()
	cancel := s.cancel
	cancel()
	return nil
}
