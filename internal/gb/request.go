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

const (
	branch = "z9hG4bK"

	letterBytes    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	contentTypeXML = "Application/MANSCDP+xml"
)

func sendDeviceInfoQuery(d model.Device) {
	document := etree.NewDocument()
	document.CreateProcInst("xml", "version=\"1.0\" encoding=\"GB2312\"")
	query := document.CreateElement("Query")
	query.CreateElement("CmdType").CreateText("DeviceInfo")
	query.CreateElement("SN").CreateText("701385")
	query.CreateElement("DeviceID").CreateText("44010200491118000001")
	document.Indent(2)
	xml, _ := document.WriteToString()
	createMessageRequest(d, xml)
}

func createMessageRequest(d model.Device, xml string) {
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
	requestBuilder.SetBody(xml)
	req, _ := requestBuilder.Build()
	log.Info("\n", req)
	if err := send(req); err != nil {
		panic(err)
	}
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
