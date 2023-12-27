package gbsip

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage/cache"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model/constant"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/parser"
	"github.com/ghettovoice/gosip/sip"
	"github.com/pkg/errors"
	"strings"
)

// Cmd SIP协议的指令结构
type cmd struct {
	server *Server
}

var c *cmd

func mustSetupCommand(server *Server) {
	c = &cmd{server: server}
}

// DeviceInfoQuery 查询设备信息
func DeviceInfoQuery(d model.Device) {
	document := etree.NewDocument()
	document.CreateProcInst("xml", "version=\"1.0\" encoding=\"GB2312\"")
	query := document.CreateElement("Query")
	query.CreateElement("CmdType").CreateText("DeviceInfo")
	query.CreateElement("SN").CreateText("701385")
	query.CreateElement("DeviceID").CreateText(d.DeviceId)
	document.Indent(2)
	body, _ := document.WriteToString()

	request := sipRequestFactory.createMessageRequest(d, body)
	log.Debugf("查询设备信息请求：\n", request)
	_, _ = c.server.sendRequest(request)
	DeviceCatalogQuery(d)
}

func DeviceCatalogQuery(device model.Device) {
	xml, err := parser.CreateQueryXML(parser.CatalogCmdType, "44010200491118000001")
	if err != nil {
		return
	}

	request := sipRequestFactory.createMessageRequest(device, xml)
	log.Debugf("发送设备目录查询信息：\n%s", request)
	_, err = c.server.sendRequest(request)
	if err != nil {
		log.Error(err)
	}
}

func DeviceBasicConfig(req *model.DeviceBasicConfigDto) error {
	xml, err := parser.CreateControlXml(parser.DeviceConfig, req.DeviceId, parser.WithBasicParams(req.Name, req.Expiration, req.HeartBeatInterval, req.HeartBeatCount))
	if err != nil {
		return errors.Wrap(err, "创建设备配置请求失败")
	}
	request := sipRequestFactory.createMessageRequest(req.Device, xml)
	log.Debugf("查询设备基本配置请求：\n%s", request)
	tx, err := c.server.sendRequest(request)
	if err != nil {
		log.Error(err)
		return err
	}
	response := getResponse(tx)
	if response.IsSuccess() {

	}
	return nil
}

func DeviceBasicConfigQuery(d model.Device) error {
	xml, err := parser.CreateQueryXML(parser.ConfigDownloadCmdType, d.DeviceId, parser.WithCustomKV("ConfigType", "BasicParam"))
	if err != nil {
		return errors.Wrap(err, "创建查询设备配置请求失败")
	}
	request := sipRequestFactory.createMessageRequest(d, xml)
	_, err = c.server.sendRequest(request)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func DeviceStatusQuery(d model.Device) error {
	xml, err := parser.CreateQueryXML(parser.DeviceStatusCmdType, d.DeviceId)
	if err != nil {
		return errors.Wrap(err, "创建查询设备状态请求失败")
	}
	request := sipRequestFactory.createMessageRequest(d, xml)
	_, err = c.server.sendRequest(request)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func AlarmSubscribe(device model.Device) error {
	xml, err := parser.CreateQueryXML(parser.AlarmCmdType, device.DeviceId, parser.WithAlarmQuery())
	if err != nil {
		return errors.Wrap(err, "创建报警订阅请求body失败")
	}
	request, err := sipRequestFactory.createSubscribeRequest(device, xml, "presence")
	if err != nil {
		return errors.Wrap(err, "创建报警订阅请求失败")
	}
	log.Debug("报警订阅请求：\n", request)
	tx, err := c.server.sendRequest(request)
	if err != nil {
		log.Error(err)
		return errors.Wrap(err, "发送请求失败")
	}
	response := getResponse(tx)
	if response == nil || !response.IsSuccess() {
		log.Error("发送请求失败")
		return errors.New("发送请求失败2")
	}
	log.Debug("报警订阅请求响应:\n", response)
	return nil
}

func CatalogSubscribe(device model.Device) error {
	xml, err := parser.CreateQueryXML(parser.CatalogCmdType, device.DeviceId)
	if err != nil {
		return errors.Wrap(err, "创建目录请订阅求body失败")
	}
	request, err := sipRequestFactory.createSubscribeRequest(device, xml, "Catalog")
	if err != nil {
		return errors.Wrap(err, "创建目录订阅请求失败")
	}
	tx, err := c.server.sendRequest(request)
	if err != nil {
		log.Error(err)
		return errors.Wrap(err, "发送目录订阅请求失败")
	}
	response := getResponse(tx)
	if response == nil || !response.IsSuccess() {
		log.Error("接收目录订阅消息确认超时")
		return errors.New("接收目录订阅消息确认超时")
	}
	return nil
}

func MobilePositionSubscribe(device model.Device) error {
	xml, err := parser.CreateQueryXML(parser.MobilePositionCmdType, device.DeviceId, parser.WithCustomKV("Interval", "5"))
	if err != nil {
		return errors.Wrap(err, "创建设备移动位置订阅请求body失败")
	}
	request, err := sipRequestFactory.createSubscribeRequest(device, xml, "presence")
	if err != nil {
		return errors.Wrap(err, "创建设备移动位置订阅请求失败")
	}
	tx, err := c.server.sendRequest(request)
	if err != nil {
		log.Error(err)
		return errors.Wrap(err, "发送设备移动位置订阅请求失败")
	}
	response := getResponse(tx)
	if response == nil || !response.IsSuccess() {
		log.Error("接收设备移动位置订阅确认超时")
		return errors.New("接收设备移动位置订阅确认超时")
	}
	return nil
}

func Play(device model.Device, detail model.MediaDetail, streamId, ssrc string, channelId string, rtpPort int) (model.StreamInfo, error) {
	log.Debugf("点播开始，流id: %c, 设备ip: %c, SSRC: %c, rtp端口: %d\n", streamId, device.Ip, ssrc, rtpPort)
	request := sipRequestFactory.createInviteRequest(device, detail, channelId, ssrc, rtpPort)
	log.Debugf("发送invite请求：\n%s", request)
	tx, err := c.server.sendRequest(request)
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
	err = c.server.s.Send(ackRequest)
	if err != nil {
		log.Errorf("发送ack失败", err)
		return model.StreamInfo{}, errors.WithMessage(err, "send play SipOption ack request fail")
	}

	// save stream info and SipOption transaction to cache
	info := model.MustNewStreamInfo(detail.ID, detail.Ip, streamId, ssrc)
	saveStreamInfo(info)

	callId, fromTag, toTag, branch, err := getRequestTxField(request, resp)
	if err != nil {
		return model.StreamInfo{}, err
	}
	streamSessionManage.saveStreamSession(device.DeviceId, channelId, ssrc, callId, fromTag, toTag, branch)

	return info, nil
}

func StopPlay(streamId, channelId string, device model.Device) error {
	// delete stream info in cache
	key := fmt.Sprintf("%s:%s", constant.StreamInfoPrefix, streamId)
	err := cache.Del(key)
	if err != nil {
		return err
	}

	// get SipOption tx in cache
	txInfo, err := streamSessionManage.getTx(device.DeviceId, channelId)
	if err != nil {
		return err
	}

	// generate bye request send to device
	byeRequest, err := sipRequestFactory.createByeRequest(channelId, device, txInfo)
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
	tx, err := c.server.sendRequest(byeRequest)
	if err != nil {
		log.Error("发送请求发生错误,", err)
	}

	response := getResponse(tx)

	if response == nil {
		log.Error("response is nil")
	}
	return nil
}

func Register(devies []model.Device) {
	request := sipRequestFactory.createRegisterRequest(devies)
	log.Infof("生成的register请求：\n%s", request)
}

// save stream info to cache
func saveStreamInfo(info model.StreamInfo) {
	key := fmt.Sprintf("%s:%s", constant.StreamInfoPrefix, info.Stream)
	cache.Set(key, info)
}

// get SipOption tx info by SipOption request and response
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

func ControlPTZ(d model.Device, channelId, command string, params1, params2, combineCode int) error {

	cmdStr, err := createPTZCode(command, params1, params2, combineCode)
	if err != nil {
		log.Error(err)
		return err
	}

	xml, err := parser.CreateControlXml(parser.DeviceControl, channelId, parser.WithPTZCmd(cmdStr))
	if err != nil {
		log.Error(err)
		return err
	}

	request := sipRequestFactory.createMessageRequest(d, xml)
	log.Info(request)
	_, err = c.server.sendRequest(request)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// 创建PTZ指令
// 根据gb28181协议的标准，前端指令中一共包含4个字节
func createPTZCode(command string, params1, params2, combineCode int) (string, error) {
	var ptz strings.Builder
	// gb28181协议中控制指令中的前三个字节
	// 字节1是A5，字节2是组合码，高4位由版本信息组成，版本信息为0H；低四位是校验位，校验位=(字节1的高4位+字节1的低四位+字节2的高四位) % 16
	// 所以校验码 = (0xa + 0x5 + 0) % 16 = (1010 + 0101 + 0) % 16 = 15 % 16 = 15；十进制数15转十六进制= F
	// 所以字节2 = 0F
	// 字节3是地址的低8位，这里直接设置为01
	ptz.WriteString("A50F01")
	var cmd int

	// 指令码以一个字节来表示
	// 0000 0000，高位的前两个bit不做表示
	// 所以有作用的也就是后6个bit，从高到低，这些bit分别控制云台的镜头缩小、镜头放大、上、下、左、右
	// 如果有做对应的操作，就将对应的bit位置1
	switch command {
	case "right":
		// 0000 0001
		cmd = 1
	case "left":
		// 0000 0010
		cmd = 2
	case "down":
		// 0000 0100
		cmd = 4
	case "up":
		// 0000 1000
		cmd = 8
	case "downright":
		// 0000 0101
		cmd = 5
	case "downleft":
		// 0000 0110
		cmd = 6
	case "upright":
		// 0000 1001
		cmd = 9
	case "upleft":
		// 0000 1010
		cmd = 10
	case "zoomin":
		// 0001 0000
		cmd = 16
	case "zoomout":
		// 0010 0000
		cmd = 32
	case "stop":
		cmd = 0
	default:
		return "", errors.New("不合规的控制字符串")
	}

	// 根据gb标准，字节4用于表示云台的镜头缩小、镜头放大、上、下、左、右，写入指令码的16进制数
	ptz.WriteString(fmt.Sprintf("%02X", cmd))

	log.Debug("合并字节4之后：" + ptz.String())

	// 根据gb标准，字节5用于表示水平控制速度，写入水平控制方向速度的十六进制数
	ptz.WriteString(fmt.Sprintf("%02X", params1))

	// 根据gb标准，字节6用于表示垂直控制速度，写入垂直控制方向速度的十六进制数
	ptz.WriteString(fmt.Sprintf("%02X", params2))

	// 最后字节7的高4位用于表示变倍控制速度，后4位不关注
	// 所以这里直接与0xF0做与操作，保留前4位，后4为置0
	c := combineCode & 0xF0
	ptz.WriteString(fmt.Sprintf("%02X", c))

	// 字节8用于校验位，根据gb标准，校验位=(字节1+字节2+字节3+字节4+字节5+字节6+字节7) % 256
	checkCode := (0xA5 + 0x0F + 0x01 + cmd + params1 + params2 + c) % 0x100
	ptz.WriteString(fmt.Sprintf("%02X", checkCode))
	log.Debug("最终生成的PTZCmd: " + ptz.String())
	return ptz.String(), nil
}
