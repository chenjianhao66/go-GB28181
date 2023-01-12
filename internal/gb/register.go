package gb

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/parser"
	"github.com/chenjianhao66/go-GB28181/internal/service"
	"github.com/ghettovoice/gosip/sip"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

const (
	letterBytes      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultAlgorithm = "MD5"
)

const (
	WWWHeader     = "WWW-Authenticate"
	ExpiresHeader = "Expires"
)

func RegisterHandler(req sip.Request, tx sip.ServerTransaction) {
	log.Infof("收到来自%s 的请求: \n%s\n", req.Source(), req.String())
	// 判断是否存在 Authorization 字段
	if headers := req.GetHeaders("Authorization"); len(headers) > 0 {
		// 存在 Authorization 头部字段
		//authHeader := headers[0].(*sip.GenericHeader)
		fromRequest, ok := parser.ParserDeviceFromRequest(req)
		if !ok {
			return
		}
		log.Infof("fromRequest: %+v\n", fromRequest)
		offlineFlag := false
		device, ok := service.Device().GetByDeviceId(fromRequest.DeviceId)

		if !ok {
			logrus.Debugln("not found from device from database")
			device = fromRequest
		}

		h := req.GetHeaders(ExpiresHeader)
		if len(h) != 1 {
			logrus.Debugln("not found expires header from request", req)
			return
		}
		expires := h[0].(*sip.Expires)
		// 如果v=0，则代表该请求是注销请求
		if expires.Equals(new(sip.Expires)) {
			logrus.Info("expires值为0,该请求是注销请求")
			offlineFlag = true
		}
		device.Expires = expires.Value()
		logrus.Debugf("设备信息:  %v\n", device)
		// 发送OK信息
		_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "ok", ""))

		if offlineFlag {
			// 注销请求
			_ = service.Device().Offline(device)
		} else {
			// 注册请求
			//device.RegisterTime = time.Now()
			if err := service.Device().Online(device); err != nil {
				logrus.Errorf("设备上线失败请检查,%s", err)
			}
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
			RandString(32),
		),
	}
	response.AppendHeader(wwwHeader)
	log.Infof("没有Authorization头部信息，生成WWW-Authenticate头部返回：\n%s\n", response)
	_ = tx.Respond(response)
}

// RandString https://github.com/kpbird/golang_random_string
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	output := make([]byte, n)
	randomness := make([]byte, n)

	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}
	l := len(letterBytes)

	for pos := range output {
		random := randomness[pos]
		randomPos := random % uint8(l)
		output[pos] = letterBytes[randomPos]
	}

	return string(output)
}
