package v1

import (
	srv "github.com/chenjianhao66/go-GB28181/internal/service/v1"
	"github.com/chenjianhao66/go-GB28181/internal/store"
)

// DeviceController 设备控制器
type DeviceController struct {
	srv srv.Service
}

// NewDeviceController 新建设备控制器
func NewDeviceController(store store.Factory) *DeviceController {
	return &DeviceController{
		srv: srv.NewService(store),
	}
}
