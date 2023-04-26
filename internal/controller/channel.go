package controller

import (
	srv "github.com/chenjianhao66/go-GB28181/internal/service"
	"github.com/chenjianhao66/go-GB28181/internal/storage"
	"github.com/gin-gonic/gin"
)

type ChannelController struct {
	srv srv.Service
}

func NewChannelController(factory storage.Factory) *ChannelController {
	return &ChannelController{
		srv: srv.NewService(factory),
	}
}

// List 返回一个设备下的所有通道
//
//	@Summary      返回一个设备下的所有通道信息
//	@Description  给定一个设备id，返回该设备下的所有通道信息
//	@Tags         设备通道
//	@Param        device    path     string  true  "设备id"
//	@Success      200  {array}   model.Channel
//	@Router       /channel/list/{device} [get]
func (c *ChannelController) List(ctx *gin.Context) {
	d := ctx.Param("device")
	if d == "" {
		newResponse(ctx).fail("device 参数是必须的")
		return
	}

	list, err := c.srv.Channel().List(d)
	if err != nil {
		newResponse(ctx).fail("查询数据库出错")
		return
	}
	newResponse(ctx).successWithAny(list)
}
