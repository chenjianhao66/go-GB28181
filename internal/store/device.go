package store

import "github.com/chenjianhao66/go-GB28181/internal/model"

// DeviceStore defines device storage interface
type DeviceStore interface {
	Store[model.Device]
	GetByDeviceId(deviceId string) (model.Device, bool)
	Keepalive(id uint) error
}
