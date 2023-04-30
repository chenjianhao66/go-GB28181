package gbsip

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/ghettovoice/gosip"
	l "github.com/ghettovoice/gosip/log"
	"github.com/ghettovoice/gosip/sip"
)

type Server struct {
	host    string
	network string
	s       gosip.Server
	c       *SipConfig
}

type RequestHandlerMap map[sip.RequestMethod]func(req sip.Request, tx sip.ServerTransaction)

type SipConfig struct {
	SipOption   *option.SIPOptions
	MysqlOption *option.MySQLOptions
	HandlerMap  RequestHandlerMap
}

func NewServer(c *SipConfig) *Server {
	s := &Server{
		host: c.SipOption.Ip + ":" + c.SipOption.Port,
		s: gosip.NewServer(
			gosip.ServerConfig{
				UserAgent: c.SipOption.UserAgent,
			},
			nil,
			nil,
			l.NewDefaultLogrusLogger(),
		),
		c: c,
	}
	s.registerHandler()
	mustSetupCommand(s)
	return s
}

func (s *Server) ListenTCP() error {
	return s.s.Listen("tcp", s.host, nil)
}

func (s *Server) ListenUDP() error {
	return s.s.Listen("udp", s.host, nil)
}

func (s *Server) Shutdown() error {
	s.s.Shutdown()
	log.Info("gb server shutdown...")
	return nil
}

func (s *Server) sendRequest(request sip.Request) (sip.ClientTransaction, error) {
	return s.s.Request(request)
}

func (s *Server) registerHandler() {
	for method, f := range s.c.HandlerMap {
		_ = s.s.OnRequest(method, f)
	}
}

func (s *Server) Register(method sip.RequestMethod, handler gosip.RequestHandler) {
	_ = s.s.OnRequest(method, handler)
}
