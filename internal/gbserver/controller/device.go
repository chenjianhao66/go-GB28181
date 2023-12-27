package controller

import (
	"fmt"
	srv "github.com/chenjianhao66/go-GB28181/internal/gbserver/service"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/syn"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"time"
)

var (
	errDeviceBasicConfig             = errors.New("修改设备基本配置失败")
	errDeviceBasicConfigQuery        = errors.New("获取设备基本配置失败")
	errDeviceBasicConfigQueryTimeOut = errors.New("获取设备基本配置超时")

	errDeviceStatusQuery        = errors.New("获取设备状态失败")
	errDeviceStatusQueryTimeOut = errors.New("获取设备状态失败")
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

	entity := syn.NewDelayTask(fmt.Sprintf("%s_%s", syn.KeyControlDeviceConfigQuery, deviceId), 3*time.Second)
	data, err := entity.Wait()
	if err != nil {
		if errors.Is(err, syn.ErrTimeOut) {
			newResponse(ctx).fail(errDeviceBasicConfigQueryTimeOut.Error())
			return
		}
		log.Error(err)
		newResponse(ctx).fail(errDeviceBasicConfigQuery.Error())
		return
	}
	newResponse(ctx).successWithAny(data)
}

func (d *DeviceController) StatusQuery(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")
	device, ok := d.srv.Devices().GetByDeviceId(deviceId)
	if !ok {
		newResponse(ctx).fail(errDeviceNotFound.Error())
		return
	}

	entity := syn.NewDelayTask(fmt.Sprintf("%s_%s", syn.KeyQueryDeviceStatus, deviceId), 3*time.Second)
	if err := gbsip.DeviceStatusQuery(device); err != nil {
		log.Error(err)
		newResponse(ctx).fail(errDeviceStatusQuery.Error())
		return
	}
	data, err := entity.Wait()
	if err != nil {
		if errors.Is(err, syn.ErrTimeOut) {
			newResponse(ctx).fail(errDeviceStatusQueryTimeOut.Error())
			return
		}
		log.Error(err)
		newResponse(ctx).fail(errDeviceStatusQuery.Error())
		return
	}

	newResponse(ctx).successWithAny(data)
}

func (d *DeviceController) AlarmSubscribe(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")
	device, ok := d.srv.Devices().GetByDeviceId(deviceId)
	if !ok {
		newResponse(ctx).fail(errDeviceNotFound.Error())
		return
	}
	if err := gbsip.AlarmSubscribe(device); err != nil {
		newResponse(ctx).fail("订阅失败")
		return
	}
	newResponse(ctx).success()
}

func (d *DeviceController) CatalogSubscribe(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")
	device, ok := d.srv.Devices().GetByDeviceId(deviceId)
	if !ok {
		newResponse(ctx).fail(errDeviceNotFound.Error())
		return
	}
	if err := gbsip.CatalogSubscribe(device); err != nil {
		newResponse(ctx).fail("订阅失败")
		return
	}
	newResponse(ctx).success()
}

func (d *DeviceController) MobilePositionSubscribe(ctx *gin.Context) {
	deviceId := ctx.Param("deviceId")
	device, ok := d.srv.Devices().GetByDeviceId(deviceId)
	if !ok {
		newResponse(ctx).fail(errDeviceNotFound.Error())
		return
	}
	if err := gbsip.MobilePositionSubscribe(device); err != nil {
		newResponse(ctx).fail("订阅失败")
		return
	}
	newResponse(ctx).success()
}
