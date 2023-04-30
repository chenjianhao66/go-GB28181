package controller

import (
	srv "github.com/chenjianhao66/go-GB28181/internal/gbserver/service"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
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

// Play 播放视频
//
// @Summary      播放设备的通道视频
// @Description  根据设备id、通道id去播放视频
// @Tags         播放
// @Produce      json
// @Param       deviceId	path	string	true	"设备id"
// @Param       channelId	path	string	true	"通道id"
// @Success      200  {object}  model.StreamInfo
// @Router       /play/start/{deviceId}/{channelId} [post]
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
