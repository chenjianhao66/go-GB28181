package gb

import (
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage/sqlite"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/ghettovoice/gosip/sip"
)

type Server struct {
	server *gbsip.Server
}

type SipConfig struct {
	SipOption   *option.SIPOptions
	MysqlOption *option.MySQLOptions
}

func NewServer(c *SipConfig) *Server {
	s := &Server{
		gbsip.NewServer(
			&gbsip.SipConfig{
				SipOption:   c.SipOption,
				MysqlOption: c.MysqlOption,
				HandlerMap:  createHandlerMap(),
			}),
	}
	// 使用sqlite
	//storage.s = mysql.GetMySQLFactory()
	storage.s = sqlite.GetSqliteFactory()
	return s
}

func (s *Server) ListenTCP() error {
	return s.server.ListenTCP()
}

func (s *Server) ListenUDP() error {
	return s.server.ListenUDP()
}

func (s *Server) Close() error {
	_ = s.server.Shutdown()
	log.Info("gb server shutdown...")
	return nil
}

func createHandlerMap() gbsip.RequestHandlerMap {
	m := make(map[sip.RequestMethod]func(req sip.Request, tx sip.ServerTransaction))
	m[sip.REGISTER] = RegisterHandler
	m[sip.MESSAGE] = MessageHandler
	return m
}
