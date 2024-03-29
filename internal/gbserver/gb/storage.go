package gb

import (
	st "github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/cron"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"time"
)

type data struct {
	s st.Factory
}

var storage = new(data)

// 设备离线
func (d *data) deviceOffline(device model.Device) error {
	log.Infof("%s设备离线,设备信息：%+v", device.DeviceId, device)
	device.Offline = 0
	err := d.s.Devices().Update(device)
	if err != nil {
		log.Errorf("设备离线发生错误，请检查。%s", err)
		return err
	}
	return nil
}

// 设备上线
func (d *data) deviceOnline(device model.Device) error {
	var err error
	if device.RegisterTime.Equal(time.Time{}) {
		log.Infof("%s设备第一次注册，发送设备查询请求", device.DeviceId)
		device.RegisterTime = time.Now()
		device.Keepalive = time.Now()
		device.Offline = 1
		err = d.s.Devices().Save(device)
	} else {
		log.Infof("%s设备离线状态下重新上线，", device.DeviceId)
		device.Offline = 1
		err = d.s.Devices().Update(device)
	}

	err = cron.StartTask(device.DeviceId, cron.TaskKeepLive, 10*time.Second, func() {
		device.Offline = 0
		d.s.Devices().Save(device)
	})

	return err
}

func (d *data) deviceKeepalive(deviceId uint) error {
	return d.s.Devices().Keepalive(deviceId)
}

// 根据设备id获取设备
func (d *data) getDeviceById(deviceId string) (model.Device, bool) {
	return d.s.Devices().GetByDeviceId(deviceId)
}

func (d *data) updateDeviceInfo(device model.Device) error {
	return d.s.Devices().UpdateDeviceInfo(device)
}

func (d *data) updateDeviceBasicConfig(device gbsip.DeviceBasicConfigResp) error {
	dev := model.Device{
		DeviceId:          device.DeviceID.DeviceID,
		Name:              device.BasicParam.Name,
		Expires:           device.BasicParam.Expiration,
		HeartBeatInterval: device.BasicParam.HeartBeatInterval,
		HeartBeatCount:    device.BasicParam.HeartBeatCount,
	}
	return d.s.Devices().UpdateBasicConfig(dev)
}

func (d *data) syncChannel(c DeviceCatalogResponse) {
	var channels []model.Channel
	for _, item := range c.DeviceList.Items {
		c := item.ConvertToChannel()
		channels = append(channels, c)
	}
	_ = d.s.Channel().SaveBatch(channels, c.DeviceID.DeviceID)
}
