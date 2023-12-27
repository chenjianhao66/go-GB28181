package mysql

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"gorm.io/gorm"
)

type channelStorage struct {
	db *gorm.DB
}

func newChannelStorage(ds *datastore) *channelStorage {
	return &channelStorage{db: ds.db}
}

func (c channelStorage) SaveBatch(channels []model.Channel, deviceId string) error {
	var deviceIds []string
	for _, c := range channels {
		deviceIds = append(deviceIds, c.DeviceId)
	}
	err := c.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("parentId = ?", deviceId).Delete(&model.Channel{}).Error; err != nil {
			return err
		}

		if err := tx.Create(channels).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c channelStorage) List(deviceId string) ([]model.Channel, error) {
	var list []model.Channel
	err := c.db.Model(&model.Channel{}).Where("parentId = ?", deviceId).Find(&list).Error
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return list, nil
}
