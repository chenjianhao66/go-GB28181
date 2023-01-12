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

const defaultConfigName = "application"

func NewServer() *Server {
	loadConfig()
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

func loadConfig() {
	viper.AddConfigPath("config/")
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("load config fail,please check your config file whether in config/ in the directory")
	}
}
