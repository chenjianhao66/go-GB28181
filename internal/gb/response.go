package gb

import (
	"encoding/xml"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
)

const (
	resultOK = "OK"
)

type deviceInfo struct {
	CmdType      string `xml:"CmdType"`
	SN           string `xml:"SN"`
	DeviceID     string `xml:"DeviceID"`
	Result       string `xml:"Result"`
	DeviceName   string `xml:"DeviceName"`
	Manufacturer string `xml:"Manufacturer"`
	Model        string `xml:"Model"`
	Firmware     string `xml:"Firmware"`
}

func deviceInfoHandler(req sip.Request, tx sip.ServerTransaction) {
	d := &deviceInfo{}
	if err := xml.Unmarshal([]byte(req.Body()), d); err != nil {
		log.Error("解析deviceInfo响应包出错", err)
		return
	}

	if d.Result != resultOK {
		log.Errorf("查询设备信息请求结果为%s，请检查", d.Result)
		return
	}

	dev := model.Device{
		Name:         d.DeviceName,
		Manufacturer: d.Manufacturer,
		Model:        d.Model,
		Firmware:     d.Firmware,
		DeviceId:     d.DeviceID,
	}

	if err := storage.updateDeviceInfo(dev); err != nil {
		log.Error("更新设备信息失败", err)
		return
	}
	response := sip.NewResponseFromRequest("", req, sip.StatusCode(http.StatusOK), "OK", "")
	err := tx.Respond(response)
	//_, err := s.s.RespondOnRequest(req, http.StatusOK, http.StatusText(http.StatusOK), "", nil)
	if err != nil {
		panic(err)
	}
}
