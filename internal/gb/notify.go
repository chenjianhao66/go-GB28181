package gb

import (
	"encoding/xml"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/parser"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
)

// xml解析心跳包结构
type keepalive struct {
	CmdType  string `xml:"CmdType"`
	SN       int    `xml:"SN"`
	DeviceID string `xml:"DeviceID"`
	Status   string `xml:"Status"`
	Info     string `xml:"Info"`
}

func keepaliveHandler(req sip.Request, tx sip.ServerTransaction) {
	keepalive := &keepalive{}
	if err := xml.Unmarshal([]byte(req.Body()), keepalive); err != nil {
		log.Debugf("keepalive 消息解析xml失败：%s", err)
		return
	}
	device, ok := parser.DeviceFromRequest(req)
	if !ok {
		return
	}
	device, ok = storage.getDeviceById(device.DeviceId)
	if !ok {
		log.Debugf("{%s}设备不存在", device.DeviceId)
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusNotFound, "device "+device.DeviceId+"not found", ""))
		return
	}

	// 更新心跳时间
	if err := storage.deviceKeepalive(device.ID); err != nil {
		log.Debugf("{%d,%s}更新心跳失败：%s", device.ID, device.DeviceId, err)
	}
	_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, http.StatusText(http.StatusOK), ""))
}
