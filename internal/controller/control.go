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
