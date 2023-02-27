package mysql

import (
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"gorm.io/gorm"
)

type channelStorage struct {
	db *gorm.DB
}

func newChannelStorage(ds *datastore) *channelStorage {
	return &channelStorage{db: ds.db}
}

func (c channelStorage) SaveBatch(channels []model.Channel) error {
	var deviceIds []string
	for _, c := range channels {
		deviceIds = append(deviceIds, c.DeviceId)
	}
	err := c.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("deviceId = ?", deviceIds).Delete(&model.Channel{}).Error; err != nil {
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
