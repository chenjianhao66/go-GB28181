package service

import (
	"github.com/chenjianhao66/go-GB28181/internal/storage"
	"github.com/chenjianhao66/go-GB28181/internal/storage/mysql"
)

type Service interface {
	Devices() IDevice
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

func (s *service) Media() IMedia {
	return Media()
}

func init() {
	factory := mysql.GetMySQLFactory()
	dService.store = factory
	mService.store = factory
}
