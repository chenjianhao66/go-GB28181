package gb

import (
	"github.com/beevik/etree"
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/ghettovoice/gosip/sip"
	"math/rand"
	"strconv"
	"time"
)

type (
	Provider struct{}
	Sender   struct{}
)

const (
	branch = "z9hG4bK"

	letterBytes    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	contentTypeXML = "Application/MANSCDP+xml"
)

var (
	RequestProvider Provider
	SipSender       Sender
)

// 发送sip协议请求
func (s Sender) transmitRequest(req sip.Request) {
	log.Info("发送SIP Request消息，Method为: ", req.Method())
	if err := send(req); err != nil {
		panic(err)
	}
}

// CreateMessageRequest 创建Message类型请求
func (p Provider) CreateMessageRequest(d model.Device, body string) sip.Request {
	requestBuilder := sip.NewRequestBuilder()
	requestBuilder.SetFrom(newFromAddress(newParams(map[string]string{"tag": randString(32)})))

	to := newTo(d.DeviceId, d.Ip)
	requestBuilder.SetTo(to)
	requestBuilder.SetRecipient(to.Uri)
	requestBuilder.AddVia(newVia(d.Transport))
	contentType := sip.ContentType(contentTypeXML)
	requestBuilder.SetContentType(&contentType)
	requestBuilder.SetMethod(sip.MESSAGE)
	userAgent := sip.UserAgentHeader("go-gb")
	requestBuilder.SetUserAgent(&userAgent)
	requestBuilder.SetBody(body)
	req, _ := requestBuilder.Build()
	return req
}

// 查询设备信息
func deviceInfoQuery(d model.Device) {
	document := etree.NewDocument()
	document.CreateProcInst("xml", "version=\"1.0\" encoding=\"GB2312\"")
	query := document.CreateElement("Query")
	query.CreateElement("CmdType").CreateText("DeviceInfo")
	query.CreateElement("SN").CreateText("701385")
	query.CreateElement("DeviceID").CreateText(d.DeviceId)
	document.Indent(2)
	xml, _ := document.WriteToString()
	request := RequestProvider.CreateMessageRequest(d, xml)
	SipSender.transmitRequest(request)
}

// 从自身SIP服务获取地址返回FromHeader
func newFromAddress(params sip.Params) *sip.Address {
	return &sip.Address{
		Uri: &sip.SipUri{
			FUser: sip.String{Str: config.SIPUser()},
			FHost: config.SIPDomain(),
		},
		Params: params,
	}
}

func newTo(user, host string) *sip.Address {
	return &sip.Address{
		Uri: &sip.SipUri{
			FUser: sip.String{Str: user},
			FHost: host,
		},
	}
}

func newParams(m map[string]string) sip.Params {
	params := sip.NewParams()
	for k, v := range m {
		params.Add(k, sip.String{Str: v})
	}
	return params
}

func newVia(transport string) *sip.ViaHop {
	port, err := strconv.ParseInt(config.SIPPort(), 10, 64)
	if err != nil {
		log.Error("解析Via头部端口失败", err)
	}
	p := sip.Port(port)

	params := newParams(map[string]string{
		"branch": branch + strconv.Itoa(int(time.Now().UnixMilli())),
	})

	return &sip.ViaHop{
		ProtocolName:    "SIP",
		ProtocolVersion: "2.0",
		Transport:       transport,
		Host:            config.SIPAddress(),
		Port:            &p,
		Params:          params,
	}
}

// randString https://github.com/kpbird/golang_random_string
func randString(n int) string {
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
