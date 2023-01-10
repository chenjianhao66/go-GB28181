package mysql

import (
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"gorm.io/gorm"
)

type devices struct {
	db *gorm.DB
}

func newDevices(ds *datastore) *devices {
	return &devices{ds.db}
}

func (d *devices) Save(entity model.Device) error {
	return d.db.Save(entity).Error
}

func (d *devices) DeleteById(id uint) error {
	return d.db.Delete(&model.Device{}, id).Error
}

func (d *devices) Update(entity model.Device) error {
	return d.db.Save(entity).Error
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
