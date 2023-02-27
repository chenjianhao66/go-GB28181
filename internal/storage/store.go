package storage

import "github.com/chenjianhao66/go-GB28181/internal/model"

// Factory defines the factory storage interface
type Factory interface {
	Devices() DeviceStore
	Media() MediaStorage
	Channel() ChannelStore
}

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

type MediaStorage interface {
	GetMediaByID(id string) (model.MediaDetail, error)
	Save(config model.MediaDetail) error
}

type ChannelStore interface {
	SaveBatch(channels []model.Channel) error
}
