package controller

import (
	"github.com/chenjianhao66/go-GB28181/internal/log"
	srv "github.com/chenjianhao66/go-GB28181/internal/service"
	"github.com/chenjianhao66/go-GB28181/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
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

// List 返回所有设备
// @Summary      返回连接到该服务的所有设备
// @Description  返回连接到该服务的所有设备
// @Tags         device
// @Accept       json
// @Produce      json
// @Success      200  {array}  model.Device
// @Router       /device/list [get]
func (d *DeviceController) List(c *gin.Context) {
	list, err := d.srv.Devices().List()
	if err != nil {
		log.Error(err)
		c.JSON(500, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}
