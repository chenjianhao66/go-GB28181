package controller

import (
	srv "github.com/chenjianhao66/go-GB28181/internal/gbserver/service"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	errDeviceBasicConfig = errors.New("修改设备基本配置失败")
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
// @Tags         设备
// @Accept       json
// @Produce      json
// @Success      200  {array}  model.Device
// @Router       /device/list [get]
func (d *DeviceController) List(c *gin.Context) {
	list, err := d.srv.Devices().List()
	if err != nil {
		log.Error(err)
		newResponse(c).fail("")
		return
	}

	newResponse(c).successWithAny(list)
}

func (d *DeviceController) BasicParamsConfig(ctx *gin.Context) {
	cfg := &model.DeviceBasicConfigReq{}
	if err := ctx.ShouldBindJSON(cfg); err != nil {
		log.Error(err)
		newResponse(ctx).fail("请求失败")
		return
	}

	device, ok := d.srv.Devices().GetByDeviceId(cfg.DeviceId)
	if !ok {
		newResponse(ctx).fail(errDeviceNotFound.Error())
		return
	} else if device.Offline == 0 {

	}

	err := gbsip.DeviceBasicConfig(&model.DeviceBasicConfigDto{*cfg, device})
	if err != nil {
		newResponse(ctx).fail(errDeviceBasicConfig.Error())
		return
	}

	newResponse(ctx).success()
}

func (d *DeviceController) BasicParamsQuery(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")
	device, ok := d.srv.Devices().GetByDeviceId(deviceId)
	if !ok {
		newResponse(ctx).fail(errDeviceNotFound.Error())
		return
	}
	if err := gbsip.DeviceBasicConfigQuery(device); err != nil {
		log.Error(err)
		newResponse(ctx).fail("fail")
		return
	}

	newResponse(ctx).success()
}
