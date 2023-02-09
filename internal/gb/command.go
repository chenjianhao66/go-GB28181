package gb

import (
	"github.com/beevik/etree"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
)

// SIPCommand SIP协议的指令结构
type SIPCommand struct{}

var sipCommand SIPCommand

// 查询设备信息
func (s SIPCommand) deviceInfoQuery(d model.Device) {
	document := etree.NewDocument()
	document.CreateProcInst("xml", "version=\"1.0\" encoding=\"GB2312\"")
	query := document.CreateElement("Query")
	query.CreateElement("CmdType").CreateText("DeviceInfo")
	query.CreateElement("SN").CreateText("701385")
	query.CreateElement("DeviceID").CreateText(d.DeviceId)
	document.Indent(2)
	body, _ := document.WriteToString()
	request := SipFactory.CreateMessageRequest(d, body)
	log.Debugf("查询设备信息请求：\n", request)
	SipSender.transmitRequest(request, nil)
}
