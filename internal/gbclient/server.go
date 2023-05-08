package gbclient

import (
	"context"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/ghettovoice/gosip/sip"
)

type (
	server struct {
		gb  *gbsip.Server
		ctx context.Context
	}
)

func newServer(c *ctlOption) *server {
	config := &gbsip.SipConfig{
		SipOption:   c.Sip,
		MysqlOption: nil,
		HandlerMap:  createHandlerMap(),
	}
	return &server{
		gb:  gbsip.NewServer(config),
		ctx: context.Background(),
	}
}

func (s *server) run() {

}

func createHandlerMap() gbsip.RequestHandlerMap {
	m := make(map[sip.RequestMethod]func(req sip.Request, tx sip.ServerTransaction))
	m[sip.REGISTER] = registerHandler
	//m[sip.MESSAGE] = MessageHandler
	return m
}
