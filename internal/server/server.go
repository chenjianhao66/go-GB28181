package server

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/gb"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	sip       *gb.Server
	apiServer *apiServer
}

func NewServer() *Server {
	return &Server{
		sip:       gb.NewServer(),
		apiServer: newApiServer(),
	}
}

func (s *Server) Run() {
	s.apiServer.initRoute()
	var eg errgroup.Group

	eg.Go(func() error {
		sport := viper.Get("server.port")
		fmt.Println(sport)
		return s.apiServer.engine.Run("127.0.0.1:18080")
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
