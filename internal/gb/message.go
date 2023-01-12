package gb

import (
	"encoding/xml"
	"github.com/chenjianhao66/go-GB28181/internal/parser"
	"github.com/chenjianhao66/go-GB28181/internal/service"
	"github.com/ghettovoice/gosip"
	"github.com/ghettovoice/gosip/sip"
	"github.com/sirupsen/logrus"
	"net/http"
)

// messageReceive 接收到的请求数据最外层，主要用来判断数据类型
type messageReceive struct {
	CmdType string `xml:"CmdType"`
	SN      int    `xml:"SN"`
}

var (
	log            = logrus.New()
	messageHandler = map[string]gosip.RequestHandler{
		"Keepalive": KeepaliveHandler,
	}
)

func MessageHandler(req sip.Request, tx sip.ServerTransaction) {
	logrus.Infof("处理MESSAGE消息....\n%s", req)
	if l, ok := req.ContentLength(); !ok || l.Equals(0) {
		logrus.Debug("该MESSAGE消息的消息体长度为0，返回OK")
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, http.StatusText(http.StatusOK), ""))
	}
	body := req.Body()
	message := &messageReceive{}
	if err := xml.Unmarshal([]byte(body), message); err != nil {
		logrus.Debugf("解析xml配置文件错误,%s\n", err)
		return
	}
	handler := messageHandler[message.CmdType]
	handler(req, tx)
}

// xml解析心跳包结构
type keepalive struct {
	CmdType  string `xml:"CmdType"`
	SN       int    `xml:"SN"`
	DeviceID string `xml:"DeviceID"`
	Status   string `xml:"Status"`
	Info     string `xml:"Info"`
}

func KeepaliveHandler(req sip.Request, tx sip.ServerTransaction) {
	keepalive := &keepalive{}
	if err := xml.Unmarshal([]byte(req.Body()), keepalive); err != nil {
		logrus.Debugf("keepalive 消息解析xml失败：%s", err)
		return
	}
	device, ok := parser.ParserDeviceFromRequest(req)
	if !ok {
		return
	}
	device, ok = service.Device().GetByDeviceId(device.DeviceId)
	if !ok {
		logrus.Debugf("{%s}设备不存在", device.DeviceId)
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusNotFound, "device "+device.DeviceId+"not found", ""))
		return
	}

	// 更新心跳时间
	if err := service.Device().Keepalive(device.ID); err != nil {
		logrus.Debugf("{%d,%s}更新心跳失败：%s", device.ID, device.DeviceId, err)
	}
	_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, http.StatusText(http.StatusOK), ""))
}
