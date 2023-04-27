package gb

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/cron"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/parser"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
)

const (
	DefaultAlgorithm = "MD5"
	WWWHeader        = "WWW-Authenticate"
	ExpiresHeader    = "Expires"
)

func RegisterHandler(req sip.Request, tx sip.ServerTransaction) {
	log.Infof("收到来自%s的请求\n", req.Source())
	// 判断是否存在 Authorization 字段
	if headers := req.GetHeaders("Authorization"); len(headers) > 0 {
		// 存在 Authorization 头部字段
		//authHeader := headers[0].(*sip.GenericHeader)
		fromRequest, ok := parser.DeviceFromRequest(req)
		if !ok {
			return
		}
		offlineFlag := false
		device, ok := storage.getDeviceById(fromRequest.DeviceId)

		if !ok {
			log.Debug("not found from device from database")
			device = fromRequest
		}

		h := req.GetHeaders(ExpiresHeader)
		if len(h) != 1 {
			log.Error("not found expires header from request", req)
			return
		}
		expires := h[0].(*sip.Expires)
		// 如果v=0，则代表该请求是注销请求
		if expires.Equals(new(sip.Expires)) {
			log.Debug("expires值为0,该请求是注销请求")
			offlineFlag = true
		}
		device.Expires = expires.Value()
		log.Infof("设备信息:  %+v\n", device)
		// 发送OK信息
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "ok", ""))

		if offlineFlag {
			// 注销请求
			_ = storage.deviceOffline(device)
			if err := cron.StopTask(device.DeviceId, cron.TaskKeepLive); err != nil {
				log.Errorf("停止心跳检测任务失败: %s", device.DeviceId)
			}
		} else {
			// 注册请求
			if err := storage.deviceOnline(device); err != nil {
				log.Errorf("设备上线失败请检查,%s", err)
			}
			go SipCommand.deviceInfoQuery(device)
		}
		return
	}

	// 没有存在 Authorization 头部字段
	response := sip.NewResponseFromRequest("", req, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), "")
	// 添加 WWW-Authenticate 头
	wwwHeader := &sip.GenericHeader{
		HeaderName: WWWHeader,
		Contents: fmt.Sprintf("Digest nonce=\"%s\", algorithm=%s, realm=\"%s\", qop=\"auth\"",
			"44010200491118000001",
			DefaultAlgorithm,
			randString(32),
		),
	}
	response.AppendHeader(wwwHeader)
	log.Debug("没有Authorization头部信息，生成WWW-Authenticate头部返回：")
	_ = tx.Respond(response)
}
