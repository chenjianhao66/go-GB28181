package mysql

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"gorm.io/gorm"
	"time"
)

type devices struct {
	db *gorm.DB
}

func newDevices(ds *datastore) *devices {
	return &devices{ds.db}
}

func (d *devices) Save(entity model.Device) error {
	return d.db.Save(&entity).Error
}

func (d *devices) DeleteById(id uint) error {
	return d.db.Delete(&model.Device{}, id).Error
}

func (d *devices) Update(entity model.Device) error {
	return d.db.Save(&entity).Error
}

func (d *devices) List() ([]model.Device, error) {
	var devices []model.Device
	if err := d.db.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (d *devices) GetById(id uint) (model.Device, error) {
	var device model.Device
	if err := d.db.First(&device, id).Error; err != nil {
		return model.Device{}, err
	}
	return device, nil
}

func (d *devices) GetByDeviceId(deviceId string) (model.Device, bool) {
	var device model.Device
	if d.db.Where("deviceId = ?", deviceId).Find(&device).RowsAffected == 0 {
		return device, false
	}
	return device, true
}

func (d *devices) Keepalive(id uint) error {
	dev := &model.Device{}
	dev.ID = id
	return d.db.Model(dev).Update("keepalive", time.Now()).Error
}

func (d *devices) UpdateDeviceInfo(entity model.Device) error {
	return d.db.Model(&model.Device{}).Where("deviceId = ?", entity.DeviceId).Updates(model.Device{
		Name:         entity.Name,
		Manufacturer: entity.Manufacturer,
		Model:        entity.Model,
		Firmware:     entity.Firmware,
	}).Error
}
