package service

import (
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/chenjianhao66/go-GB28181/internal/store"
	"github.com/chenjianhao66/go-GB28181/internal/store/mysql"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type IDevice interface {
	baseService[model.Device]
	Online(device model.Device) error
	Offline(device model.Device) error
	GetByDeviceId(deviceId string) (model.Device, bool)
	Keepalive(id uint) error
}

type deviceService struct {
	store store.Factory
}

var (
	dService *deviceService
	once     sync.Once
)

func Device() IDevice {
	once.Do(func() {
		factory, err := mysql.GetMySQLFactory()
		if err != nil {
			panic(err)
		}
		dService = &deviceService{
			store: factory,
		}
	})
	return dService
}

func (d *deviceService) Save(entity model.Device) error {
	return d.store.Devices().Save(entity)
}

func (d *deviceService) DeleteById(id uint) error {
	return d.store.Devices().DeleteById(id)
}

func (d *deviceService) Update(entity model.Device) error {
	return d.store.Devices().Update(entity)
}

func (d *deviceService) List() ([]model.Device, error) {
	return d.store.Devices().List()
}

func (d *deviceService) GetById(id uint) (model.Device, error) {
	return d.store.Devices().GetById(id)
}

func (d *deviceService) Online(device model.Device) error {
	var err error
	if device.RegisterTime.Equal(time.Time{}) {
		logrus.Infof("{%s}设备第一次注册，发送设备查询请求", device.DeviceId)
		device.RegisterTime = time.Now()
		device.Keepalive = time.Now()
		device.Offline = 1
		err = d.Save(device)
	} else {
		logrus.Infof("{%s}设备离线状态下重新上线，", device.DeviceId)
		device.Offline = 1
		err = d.Update(device)
	}
	// TODO 发起查询设备信息
	return err
}

func (d *deviceService) Offline(device model.Device) error {
	logrus.Infof("{%s}设备离线,设备信息：{%#v}", device.DeviceId, device)
	device.Offline = 0
	err := d.Update(device)
	if err != nil {
		logrus.Errorf("设备离线发生错误，请检查。%s", err)
		return err
	}
	return nil
}

func (d *deviceService) GetByDeviceId(deviceId string) (model.Device, bool) {
	return d.store.Devices().GetByDeviceId(deviceId)
}

func (d *deviceService) Keepalive(id uint) error {
	return d.store.Devices().Keepalive(id)
}
