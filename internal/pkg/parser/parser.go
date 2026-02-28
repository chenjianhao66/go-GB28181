package parser

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/ghettovoice/gosip/sip"
)

func DeviceFromRequest(req sip.Request) (model.Device, bool) {
	d := model.Device{}

	from, ok := req.From()
	if !ok {
		log.Debug("从请求中无法解析from头部信息", req.String())
		return d, false
	}
	if from.Address == nil {
		log.Debug("从请求中无法解析from头address部分信息", req.String())
		return d, false
	}
	if from.Address.User() == nil {
		log.Debug("从请求中无法解析from头user部分信息", req.String())
		return d, false
	}

	d.DeviceId = from.Address.User().String()
	d.Domain = from.Address.Host()
	via, ok := req.ViaHop()
	if !ok {
		log.Debug("从请求中无法解析出via头部信息", via.String())
		return d, false
	}
	d.Ip = via.Host
	if rport, ok := via.Params.Get("rport"); ok && rport != nil {
		d.Port = rport.String()
	} else {
		d.Port = via.Port.String()
	}
	d.Transport = via.Transport
	log.Debugf("从请求中解析出的设备信息: %v\n", d)
	return d, true
}
