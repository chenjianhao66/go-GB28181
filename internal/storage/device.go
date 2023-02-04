package storage

import "github.com/chenjianhao66/go-GB28181/internal/model"

// DeviceStore defines device storage interface
type DeviceStore interface {
	Save(entity model.Device) error
	DeleteById(id uint) error
	Update(entity model.Device) error
	UpdateDeviceInfo(entity model.Device) error
	List() ([]model.Device, error)
	GetById(id uint) (model.Device, error)
	GetByDeviceId(deviceId string) (model.Device, bool)
	Keepalive(id uint) error
}
