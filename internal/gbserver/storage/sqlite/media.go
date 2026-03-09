package sqlite

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"gorm.io/gorm"
)

type mediaStorage struct {
	db *gorm.DB
}

func newMediaStorage(ds *sqliteDatastore) *mediaStorage {
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

func (m *mediaStorage) List() ([]model.MediaDetail, error) {
	var list []model.MediaDetail
	if err := m.db.Find(&list).Error; err != nil {
		log.Error("获取流媒体服务列表失败, err: ", err)
		return nil, err
	}

	return list, nil
}
