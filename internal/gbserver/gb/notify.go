package gb

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/chenjianhao66/go-GB28181/internal/gbserver/service"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/cron"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/parser"
	"github.com/ghettovoice/gosip/sip"
)

// xml解析心跳包结构
type keepalive struct {
	CmdType  string `xml:"CmdType"`
	SN       int    `xml:"SN"`
	DeviceID string `xml:"DeviceID"`
	Status   string `xml:"Status"`
	Info     string `xml:"Info"`
}

func keepaliveNotifyHandler(req sip.Request, tx sip.ServerTransaction) {
	keepalive := &keepalive{}
	body := req.Body()
	body = strings.Replace(body, "GB2312", "UTF-8", 1)
	if err := xml.Unmarshal([]byte(body), keepalive); err != nil {
		log.Debugf("keepalive 消息解析xml失败：%s", err)
		return
	}
	device, ok := parser.DeviceFromRequest(req)
	if !ok {
		return
	}
	deviceInDb, ok := storage.getDeviceById(device.DeviceId)
	if !ok {
		log.Debugf("{%s}设备不存在", device.DeviceId)
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusNotFound, "device "+device.DeviceId+"not found", ""))
		return
	}

	if deviceInDb.Ip != device.Ip || deviceInDb.Port != device.Port {
		log.Infof("设备 [%s] 地址发生变化， 新地址和端口为：[%s]:[%s]", deviceInDb.DeviceId, device.Ip, device.Port)
		deviceInDb.Ip = device.Ip
		deviceInDb.Port = device.Port
	}
	deviceInDb.Keepalive = time.Now()
	deviceInDb.Offline = 1

	if err := storage.save(deviceInDb); err != nil {
		log.Debugf("{%d,%s}更新心跳失败：%v", deviceInDb.ID, deviceInDb.DeviceId, err.Error())
	}

	_ = cron.StopTask(device.DeviceId, cron.TaskKeepLive)
	interval := time.Duration(device.HeartBeatInterval) * time.Second * 3
	if interval == 0 {
		interval = time.Minute * 3
	}
	if err := cron.StartTask(device.DeviceId, cron.TaskKeepLive, interval, func() {
		log.Infof("[%s] 设备心跳超时, 离线", deviceInDb.DeviceId)
		deviceInDb.Offline = 0
		if err := service.Device().Offline(deviceInDb); err != nil {
			log.Errorf("[%s ]更新设备离线状态失败, err: %v", deviceInDb.DeviceId, err)
		}
	}); err != nil {
		log.Errorf("启动设备 [%s] 心跳超时任务失败, err: %v", deviceInDb.DeviceId, err)
	}

	_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, http.StatusText(http.StatusOK), ""))
}

func alarmNotifyHandler(req sip.Request, tx sip.ServerTransaction) {
	// 使用 gbsip.AlarmNotify 结构体解析
	// 自行扩展

	_ = responseAck(tx, req)
}

func mobilePositionNotifyHandler(req sip.Request, tx sip.ServerTransaction) {
	// 自行扩展

	_ = responseAck(tx, req)
}

func subscribeAlarmResponseHandler(req sip.Request, tx sip.ServerTransaction) {
	r := parser.GetResultFromXML(req.Body())
	if r == "" {
		log.Error("获取不到响应信息中的Result字段")
		return
	}

	if r == "ERROR" {
		log.Error("订阅报警信息失败，请检查")
	} else {
		log.Debug("订阅报警信息成功")
	}
	_ = responseAck(tx, req)
}

func subscribeMobilePositionResponseHandler(req sip.Request, tx sip.ServerTransaction) {
	r := parser.GetResultFromXML(req.Body())
	if r == "" {
		log.Error("获取不到响应信息中的Result字段")
		return
	}

	if r == "ERROR" {
		log.Error("订阅设备移动位置信息失败，请检查")
	} else {
		log.Debug("订阅设备移动位置信息信息成功")
	}
	_ = responseAck(tx, req)
}
