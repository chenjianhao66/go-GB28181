package gb

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/chenjianhao66/go-GB28181/internal/model/constant"
	"github.com/chenjianhao66/go-GB28181/internal/parser"
	"github.com/chenjianhao66/go-GB28181/internal/storage/cache"
	"github.com/ghettovoice/gosip/sip"
	"github.com/pkg/errors"
)

// SIPCommand SIP协议的指令结构
type SIPCommand struct{}

var SipCommand SIPCommand

// DeviceInfoQuery 查询设备信息
func (c SIPCommand) deviceInfoQuery(d model.Device) {
	document := etree.NewDocument()
	document.CreateProcInst("xml", "version=\"1.0\" encoding=\"GB2312\"")
	query := document.CreateElement("Query")
	query.CreateElement("CmdType").CreateText("DeviceInfo")
	query.CreateElement("SN").CreateText("701385")
	query.CreateElement("DeviceID").CreateText(d.DeviceId)
	document.Indent(2)
	body, _ := document.WriteToString()
	request := SipFactory.createMessageRequest(d, body)
	log.Debugf("查询设备信息请求：\n", request)
	_, _ = SipSender.TransmitRequest(request)
	c.deviceCatalogQuery(d)
}

func (c SIPCommand) deviceCatalogQuery(device model.Device) {
	xml, err := parser.CreateQueryXML(parser.CatalogCmdType, "44010200491118000001")
	if err != nil {
		return
	}

	request := SipFactory.createMessageRequest(device, xml)
	log.Debugf("发送设备目录查询信息：\n%s", request)
	_, err = SipSender.TransmitRequest(request)
	if err != nil {
		log.Error(err)
	}
}

func (c SIPCommand) Play(device model.Device, detail model.MediaDetail, streamId, ssrc string, channelId string, rtpPort int) (model.StreamInfo, error) {
	log.Debugf("点播开始，流id: %c, 设备ip: %c, SSRC: %c, rtp端口: %d\n", streamId, device.Ip, ssrc, rtpPort)
	request := SipFactory.createInviteRequest(device, detail, channelId, ssrc, rtpPort)
	log.Debugf("发送invite请求：\n%s", request)
	tx, err := SipSender.TransmitRequest(request)
	if err != nil {
		return model.StreamInfo{}, err
	}

	resp := getResponse(tx)
	log.Debugf("收到invite响应：\n%s", resp)
	log.Debugf("\ntransaction key: %s", tx.Key().String())

	ackRequest := sip.NewAckRequest("", request, resp, "", nil)
	ackRequest.SetRecipient(request.Recipient())
	ackRequest.AppendHeader(&sip.ContactHeader{
		Address: request.Recipient(),
		Params:  nil,
	})

	log.Debugf("发送ack确认：%s\n", ackRequest)
	err = s.s.Send(ackRequest)
	if err != nil {
		log.Errorf("发送ack失败", err)
		return model.StreamInfo{}, errors.WithMessage(err, "send play sip ack request fail")
	}

	// save stream info and sip transaction to cache
	info := model.MustNewStreamInfo(detail.ID, detail.Ip, streamId, ssrc)
	saveStreamInfo(info)

	callId, fromTag, toTag, branch, err := getRequestTxField(request, resp)
	if err != nil {
		return model.StreamInfo{}, err
	}
	streamSessionManage.saveStreamSession(device.DeviceId, channelId, ssrc, callId, fromTag, toTag, branch)

	return info, nil
}

func (c SIPCommand) StopPlay(streamId, channelId string, device model.Device) error {
	// delete stream info in cache
	key := fmt.Sprintf("%s:%s", constant.StreamInfoPrefix, streamId)
	err := cache.Del(key)
	if err != nil {
		return err
	}

	// get sip tx in cache
	txInfo, err := streamSessionManage.getTx(device.DeviceId, channelId)
	if err != nil {
		return err
	}

	// generate bye request send to device
	byeRequest, err := SipFactory.createByeRequest(channelId, device, txInfo)
	if err != nil {
		return err
	}

	log.Debugf("创建Bye请求：\n%s", byeRequest)
	key = fmt.Sprintf("%s:%s", constant.StreamTransactionPrefix, streamId)
	err = cache.Del(key)
	if err != nil {
		return errors.WithMessage(err, "delete cache by key fail")
	}

	//err = s.s.Send(byeRequest)
	tx, err := SipSender.TransmitRequest(byeRequest)
	if err != nil {
		log.Error("发送请求发生错误,", err)
	}

	response := getResponse(tx)

	if response == nil {
		log.Error("response is nil")
	}
	return nil
}

// save stream info to cache
func saveStreamInfo(info model.StreamInfo) {
	key := fmt.Sprintf("%s:%s", constant.StreamInfoPrefix, info.Stream)
	cache.Set(key, info)
}

// get sip tx info by sip request and response
func getRequestTxField(request sip.Request, response sip.Response) (callId, fromTag, toTag, viaBranch string, err error) {
	callID, ok := request.CallID()
	if !ok {
		return "", "", "", "", errors.New("get CallId header in request fail")
	}

	fromHeader, ok := request.From()
	if !ok {
		return "", "", "", "", errors.New("get from header in request fail")
	}
	ft, ok := fromHeader.Params.Get("tag")
	if !ok {
		return "", "", "", "", errors.New("get tag field in 'from' header fail")
	}

	toHeader, ok := response.To()
	if !ok {
		return "", "", "", "", errors.New("get to header in request fail")
	}
	tg, ok := toHeader.Params.Get("tag")
	if !ok {
		return "", "", "", "", errors.New("get tag field in 'to' header fail")
	}

	viaHop, ok := request.ViaHop()
	if !ok {
		return "", "", "", "", errors.New("get via header in request fail")
	}

	branch, ok := viaHop.Params.Get("branch")
	if !ok {
		return "", "", "", "", errors.New("get branch field in 'via' header fail")
	}

	callId = callID.Value()
	fromTag = ft.String()
	toTag = tg.String()
	viaBranch = branch.String()
	return
}
