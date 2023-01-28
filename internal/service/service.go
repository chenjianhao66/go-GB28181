package service

import "github.com/chenjianhao66/go-GB28181/internal/store"

type Service interface {
	Devices() IDevice
}

type service struct {
	store store.Factory
}

// NewService 新建服务接口
func NewService(factory store.Factory) Service {
	return &service{store: factory}
}

// Devices 返回设备服务实现接口
func (s *service) Devices() IDevice {
	return Device()
}
