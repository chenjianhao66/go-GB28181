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

func (c *ChannelController) List(ctx *gin.Context) {
	d := ctx.Param("device")
	if d == "" {
		ctx.JSON(500, gin.H{
			"msg": "device 参数是必须的",
		})
		return
	}

	list, err := c.srv.Channel().List(d)
	if err != nil {
		ctx.JSON(500, gin.H{
			"msg": "查询数据库出错",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"data": list,
	})
}
