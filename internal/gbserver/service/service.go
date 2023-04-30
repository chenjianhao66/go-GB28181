package service

import (
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
)

type Service interface {
	Devices() IDevice
	Play() IPlay
	Media() IMedia
	Channel() IChannel
}

type service struct {
	store storage.Factory
}

// NewService 新建服务接口
func NewService(factory storage.Factory) Service {
	return &service{store: factory}
}

// Devices 返回设备服务实现接口
func (s *service) Devices() IDevice {
	return Device()
}

func (s *service) Play() IPlay {
	return Play()
}

func (s *service) Media() IMedia {
	return Media()
}

func (s *service) Channel() IChannel {
	return Channel()
}

func InitService(factory storage.Factory) {
	dService.store = factory
	mService.store = factory
	cService.store = factory
}
