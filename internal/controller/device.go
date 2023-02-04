package controller

import (
	srv "github.com/chenjianhao66/go-GB28181/internal/service"
	"github.com/chenjianhao66/go-GB28181/internal/storage"
	"github.com/gin-gonic/gin"
)

// DeviceController 设备控制器
type DeviceController struct {
	srv srv.Service
}

// NewDeviceController 新建设备控制器
func NewDeviceController(store storage.Factory) *DeviceController {
	return &DeviceController{
		srv: srv.NewService(store),
	}
}

func (d *DeviceController) List(c *gin.Context) {
	_ = c.Params
}
