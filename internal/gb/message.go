package gb

import (
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/parser"
	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
)

var (
	messageHandler = map[string]gosip.RequestHandler{
		// 通知
		"Notify:Keepalive": keepaliveHandler,

		// 响应
		// 查询设备信息响应
		"Response:DeviceInfo": deviceInfoHandler,
	}
)

func MessageHandler(req sip.Request, tx sip.ServerTransaction) {
	log.Debug("处理MESSAGE消息....\n%s", req)
	if l, ok := req.ContentLength(); !ok || l.Equals(0) {
		log.Debug("该MESSAGE消息的消息体长度为0，返回OK")
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, http.StatusText(http.StatusOK), ""))
	}
	body := req.Body()
	cmdType, err := parser.GetCmdTypeFromXML(body)
	log.Debug("解析出的命令：%s", cmdType)
	if err != nil {
		return
	}
	handler, ok := messageHandler[cmdType]
	if !ok {
		log.Infof("不支持的Message方法实现")
		return
	}
	handler(req, tx)
}
