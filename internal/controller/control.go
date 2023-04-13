package controller

import (
	"github.com/chenjianhao66/go-GB28181/internal/gb"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	errDataBindStructFail = errors.New("传入参数失败，请检查传参")
)

type ControlController struct {
}

func NewControlController() *ControlController {
	return &ControlController{}
}

// ControlPTZ 控制云台
//
//	@Summary      控制摄像头的云台
//	@Description  根据传入的控制方向、控制速度等参数去控制摄像头的云台
//	@Tags         设备控制
//	@Accept       json
//	@Produce      json
//	@Param        控制云台对象 body model.DeviceControl  true  "控制云台对象"
//	@Success      200  {string}   "ok"
//	@Router       /control/ptz [get]
func (c ControlController) ControlPTZ(ctx *gin.Context) {
	var data model.DeviceControl
	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Error(err)
		newResponse(ctx).fail(errDataBindStructFail.Error())
		return
	}
	err := gb.SipCommand.ControlPTZ(data.DeviceId, data.ChannelId, data.Command, data.HorizonSpeed, data.VerticalSpeed, data.ZoomSpeed)
	if err != nil {
		newResponse(ctx).fail(err.Error())
		return
	}
	newResponse(ctx).success()
}
