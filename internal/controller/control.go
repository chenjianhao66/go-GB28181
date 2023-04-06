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

//TODO 未完成

// ControlPTZ 控制云台
//
//	@Summary      控制摄像头的云台
//	@Description  根据传入的控制方向、控制速度等参数去控制摄像头的云台
//	@Tags         control
//	@Accept       json
//	@Produce      json
//	@Param        deviceControl body model.DeviceControl  true  "控制云台对象"
//	@Success      200  {string}   "model.Account"
//	@Router       /control/ptz [get]
func (c ControlController) ControlPTZ(ctx *gin.Context) {
	var data model.DeviceControl
	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Error(err)
		ctx.JSON(200, gin.H{
			"msg": errDataBindStructFail.Error(),
		})
		return
	}
	err := gb.SipCommand.ControlPTZ(data.DeviceId, data.ChannelId, data.Command, data.HorizonSpeed, data.VerticalSpeed, data.ZoomSpeed)
	if err != nil {
		ctx.JSON(500, gin.H{
			"msg": err.Error(),
		})
		return
	}
}
