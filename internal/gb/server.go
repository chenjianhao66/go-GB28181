package gb

import (
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/ghettovoice/gosip"
	l "github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
	"sync"
)

type Server struct {
	host    string
	network string
	s       gosip.Server
}

var (
	s          *Server
	serverOnce sync.Once
)

func NewServer() *Server {
	serverOnce.Do(func() {
		s = &Server{
			host: config.SIPAddress() + ":" + config.SIPPort(),
			s: gosip.NewServer(
				gosip.ServerConfig{
					UserAgent: config.SIPUserAgent(),
				},
				nil,
				nil,
				l.NewDefaultLogrusLogger(),
			),
		}
	})
	registerHandler(s)
	return s
}

func (s *Server) ListenTCP() error {
	return s.s.Listen("tcp", s.host, nil, nil)
}

func (s *Server) ListenUDP() error {
	return s.s.Listen("udp", s.host, nil, nil)
}

func (s *Server) Close() error {
	log.Info("gb sip server shutdown...")
	s.s.Shutdown()
	return nil
}

func registerHandler(s *Server) {
	_ = s.s.OnRequest(sip.REGISTER, RegisterHandler)
	_ = s.s.OnRequest(sip.MESSAGE, MessageHandler)
}

func send(message sip.Message) error {
	return s.s.Send(message)
}
