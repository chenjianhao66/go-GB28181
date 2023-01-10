package mysql

import (
	"github.com/chenjianhao66/go-GB28181/internal/store"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (d *datastore) Devices() store.DeviceStore {
	return newDevices(d)
}
