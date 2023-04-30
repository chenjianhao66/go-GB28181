package mysql

import (
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"gorm.io/gorm"
)

type mediaStorage struct {
	db *gorm.DB
}

func newMediaStorage(ds *datastore) *mediaStorage {
	return &mediaStorage{db: ds.db}
}

func (m *mediaStorage) Save(config model.MediaDetail) error {
	return m.db.Save(config).Error
}

func (m *mediaStorage) GetMediaByID(id string) (model.MediaDetail, error) {
	detail := model.MediaDetail{}
	err := m.db.Where("id = ?", id).First(&detail).Error
	return detail, err
}
