package service

import (
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/chenjianhao66/go-GB28181/internal/storage"
)

type IChannel interface {
	List(deviceId string) ([]model.Channel, error)
}

type channelService struct {
	store storage.Factory
}

var cService = new(channelService)

func Channel() IChannel {
	return cService
}

func (c channelService) List(deviceId string) ([]model.Channel, error) {
	list, err := c.store.Channel().List(deviceId)
	if err != nil {
		return nil, err
	}
	return list, nil
}
