package controller

import (
	"github.com/chenjianhao66/go-GB28181/internal/log"
	srv "github.com/chenjianhao66/go-GB28181/internal/service"
	"github.com/chenjianhao66/go-GB28181/internal/storage"
	"github.com/gin-gonic/gin"
)

// PlayController 设备控制器
type PlayController struct {
	srv srv.Service
}

// NewPlayController 新建播放控制器
func NewPlayController(store storage.Factory) *PlayController {
	return &PlayController{
		srv: srv.NewService(store),
	}
}

func (p *PlayController) Play(c *gin.Context) {
	deviceId := c.Param("deviceId")
	channelId := c.Param("channelId")
	streamInfo, err := p.srv.Play().Play(deviceId, channelId)
	if err != nil {
		log.Errorf("%+v", err)
		newResponse(c).fail(err.Error())
		return
	}
	newResponse(c).successWithAny(streamInfo)
}
