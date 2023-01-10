package gb

import (
	"github.com/ghettovoice/gosip"
	l "github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
)

type Server struct {
	host    string
	network string
	s       gosip.Server
}

func NewServer() *Server {
	s := &Server{
		host:    "0.0.0.0:5060",
		network: "udp",
		s: gosip.NewServer(
			gosip.ServerConfig{},
			nil,
			nil,
			l.NewDefaultLogrusLogger(),
		),
	}
	registerHandler(s)
	return s
}

func (s *Server) Listen() error {
	// TODO 引入配置文件动态变化
	return s.s.Listen("udp", "0.0.0.0:5060", nil, nil)
}

func registerHandler(s *Server) {
	_ = s.s.OnRequest(sip.REGISTER, RegisterHandler)
	_ = s.s.OnRequest(sip.MESSAGE, MessageHandler)
}
