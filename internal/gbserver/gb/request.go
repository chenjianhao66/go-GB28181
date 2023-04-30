package gb

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage/cache"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/ghettovoice/gosip/sip"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type (
	sipFactory struct{}
	sipSender  struct{}
)

const (
	letterBytes    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	contentTypeXML = "Application/MANSCDP+xml"
	contentTypeSDP = "APPLICATION/SDP"
)

var (
	sipRequestFactory sipFactory
	sipRequestSender  sipSender
)

// TransmitRequest 发送sip协议请求
//func (sender sipSender) TransmitRequest(req sip.Request) (sip.ClientTransaction, error) {
//	log.Info("发送SIP Request消息，Method为: ", req.Method())
//	transaction, err := s.s.Request(req)
//	return transaction, err
//}

// createMessageRequest 创建Message类型请求
func (f sipFactory) createMessageRequest(d model.Device, body string) sip.Request {
	requestBuilder := sip.NewRequestBuilder()
	requestBuilder.SetFrom(newFromAddress(newParams(map[string]string{"tag": randString(32)})))

	to := newTo(d.DeviceId, d.Ip, d.Port)
	requestBuilder.SetTo(to)
	requestBuilder.SetRecipient(to.Uri)
	requestBuilder.AddVia(newVia(d.Transport))
	contentType := sip.ContentType(contentTypeXML)
	requestBuilder.SetContentType(&contentType)
	requestBuilder.SetMethod(sip.MESSAGE)
	userAgent := sip.UserAgentHeader("go-gb")
	requestBuilder.SetUserAgent(&userAgent)
	requestBuilder.SetBody(body)

	ceq, err := cache.GetCeq()
	if err != nil {
		log.Error("get ceq in cache fail,", err)
	} else {
		requestBuilder.SetSeqNo(cast.ToUint(ceq))
	}
	req, _ := requestBuilder.Build()
	return req
}

// createInviteRequest 创建invite请求
func (f sipFactory) createInviteRequest(device model.Device, detail model.MediaDetail, channelId string, ssrc string, rtpPort int) sip.Request {
	body := createSdpInfo(detail.Ip, channelId, ssrc, rtpPort)

	requestBuilder := sip.NewRequestBuilder()
	to := newTo(channelId, device.Ip, device.Port)
	requestBuilder.SetMethod(sip.INVITE)
	requestBuilder.SetFrom(newFromAddress(newParams(map[string]string{"tag": randString(32)})))
	requestBuilder.SetTo(to)
	sipUri := &sip.SipUri{
		FUser: sip.String{Str: channelId},
		FHost: to.Uri.Host(),
	}
	requestBuilder.SetRecipient(sipUri)
	requestBuilder.AddVia(newVia("UDP"))
	requestBuilder.SetContact(newTo(config.SIPId(), config.SIPAddress(), config.SIPPort()))
	contentType := sip.ContentType(contentTypeSDP)
	requestBuilder.SetContentType(&contentType)
	requestBuilder.SetBody(body)
	ceq, err := cache.GetCeq()
	if err != nil {
		log.Error("get ceq in cache fail,", err)
	} else {
		requestBuilder.SetSeqNo(cast.ToUint(ceq))
	}
	callID := sip.CallID(fmt.Sprintf("%s", randString(32)))
	requestBuilder.SetCallID(&callID)
	header := sip.GenericHeader{
		HeaderName: "Subject",
		Contents:   fmt.Sprintf("%s:%s,%s:%d", channelId, ssrc, config.SIPId(), 0),
	}
	requestBuilder.AddHeader(&header)
	request, err := requestBuilder.Build()
	if err != nil {
		log.Error("发生错误：", err)
		return nil
	}

	return request
}

// create bye request
func (f sipFactory) createByeRequest(channelId string, device model.Device, tx SipTX) (sip.Request, error) {

	fromAddress := newFromAddress(newParams(map[string]string{"tag": tx.FromTag}))

	toAddress := newTo(channelId, device.Ip, device.Port)
	toAddress.Params = newParams(map[string]string{"tag": tx.ToTag})

	via := newVia(device.Transport)
	via.Params = newParams(map[string]string{"branch": tx.ViaBranch})

	callID := sip.CallID(tx.CallId)
	ceq, err := cache.GetCeq()
	if err != nil {
		log.Error("get ceq in cache fail,", err)
		ceq = 0
	}

	request, err := sip.NewRequestBuilder().
		SetFrom(fromAddress).
		SetTo(toAddress).
		SetMethod(sip.BYE).
		AddVia(via).
		SetContact(newTo(config.SIPId(), config.SIPAddress(), config.SIPPort())).
		SetCallID(&callID).
		SetSeqNo(cast.ToUint(ceq)).
		SetRecipient(&sip.SipUri{
			FUser: sip.String{channelId},
			FHost: device.Ip,
		}).Build()

	if err != nil {
		return nil, errors.WithMessage(err, "generate bye request fail")
	}
	return request, nil
}

// 从自身SIP服务获取地址返回FromHeader
func newFromAddress(params sip.Params) *sip.Address {
	log.Info(config.SIPId())
	return &sip.Address{
		Uri: &sip.SipUri{
			FUser: sip.String{Str: config.SIPId()},
			FHost: config.SIPDomain(),
		},
		Params: params,
	}
}

func newTo(user, host, port string) *sip.Address {
	p := sip.Port(cast.ToUint16(port))
	return &sip.Address{
		Uri: &sip.SipUri{
			FUser: sip.String{Str: user},
			FHost: host,
			FPort: &p,
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
		"branch": fmt.Sprintf("%s%d", "z9hG4bK", time.Now().UnixMilli()),
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

func getResponse(tx sip.ClientTransaction) sip.Response {
	timer := time.NewTimer(5 * time.Second)

	for {
		select {
		case resp := <-tx.Responses():
			if resp.StatusCode() == sip.StatusCode(http.StatusContinue) ||
				resp.StatusCode() == sip.StatusCode(http.StatusSwitchingProtocols) {
				continue
			}
			return resp
		case <-timer.C:
			log.Error("获取响应超时")
			return nil
		}
	}
}

func responseAck(transaction sip.ServerTransaction, request sip.Request) error {
	err := transaction.Respond(sip.NewResponseFromRequest("", request, sip.StatusCode(http.StatusOK),
		http.StatusText(http.StatusOK), ""))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
