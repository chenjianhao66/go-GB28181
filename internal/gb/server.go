package gb

import (
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/ghettovoice/gosip"
	l "github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
	"github.com/spf13/viper"
)

type Server struct {
	host    string
	network string
	s       gosip.Server
}

func NewServer() *Server {
	sipConfig := &config.SIPOptions{}
	if err := viper.UnmarshalKey("sip", sipConfig); err != nil {
		panic("load sip config fail")
	}
	s := &Server{
		host: sipConfig.Ip + ":" + sipConfig.Port,
		s: gosip.NewServer(
			gosip.ServerConfig{
				UserAgent: sipConfig.UserAgent,
			},
			nil,
			nil,
			l.NewDefaultLogrusLogger(),
		),
	}
	registerHandler(s)
	return s
}

func (s *Server) ListenTCP() error {
	return s.s.Listen("tcp", s.host, nil, nil)
}

func (s *Server) ListenUDP() error {
	return s.s.Listen("udp", s.host, nil, nil)
}

func registerHandler(s *Server) {
	_ = s.s.OnRequest(sip.REGISTER, RegisterHandler)
	_ = s.s.OnRequest(sip.MESSAGE, MessageHandler)
}
