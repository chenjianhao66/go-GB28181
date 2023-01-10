package gb

import (
	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/sip"
	"github.com/sirupsen/logrus"
)

var (
	log            = logrus.New()
	messageHandler = map[string]gosip.RequestHandler{
		"Keepalive": KeepaliveHandler,
	}
)

func MessageHandler(req sip.Request, tx sip.ServerTransaction) {

}

func KeepaliveHandler(req sip.Request, tx sip.ServerTransaction) {

}
