package gb

import (
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage/mysql"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/ghettovoice/gosip"
	l "github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
	"sync"
)

type Server struct {
	host    string
	network string
	s       gosip.Server
	c       *SipConfig
}

type SipConfig struct {
	SipOption   *option.SIPOptions
	MysqlOption *option.MySQLOptions
}

var (
	s          *Server
	serverOnce sync.Once
)

func NewServer(c *SipConfig) *Server {
	serverOnce.Do(func() {
		s = &Server{
			host: c.SipOption.Ip + ":" + c.SipOption.Port,
			s: gosip.NewServer(
				gosip.ServerConfig{
					UserAgent: c.SipOption.UserAgent,
				},
				nil,
				nil,
				l.NewDefaultLogrusLogger(),
			),
		}
		storage.s = mysql.GetMySQLFactory()
	})
	registerHandler(s)
	return s
}

func (s *Server) ListenTCP() error {
	return s.s.Listen("tcp", s.host, nil)
}

func (s *Server) ListenUDP() error {
	return s.s.Listen("udp", s.host, nil)
}

func (s *Server) Close() error {
	s.s.Shutdown()
	log.Info("gb server shutdown...")
	return nil
}

func registerHandler(s *Server) {
	_ = s.s.OnRequest(sip.REGISTER, RegisterHandler)
	_ = s.s.OnRequest(sip.MESSAGE, MessageHandler)
}

func send(message sip.Message) error {
	return s.s.Send(message)
}
