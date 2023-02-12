package controller

import (
	srv "github.com/chenjianhao66/go-GB28181/internal/service"
	"github.com/chenjianhao66/go-GB28181/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, streamInfo)
}
