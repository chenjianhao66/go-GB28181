package controller

import (
	"strings"
	"time"

	"github.com/chenjianhao66/go-GB28181/internal/gbserver/service"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/cron"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/gin-gonic/gin"
)

type MediaHookController struct{}

func NewMediaHookController() MediaHookController {
	return MediaHookController{}
}

// OnServerStarted 服务器启动事件，可以用于监听服务器崩溃重启；此事件对回复不敏感。
func (m MediaHookController) OnServerStarted(c *gin.Context) {
	log.Info("zlm 上线...")
	conf := model.MediaConfig{}

	if err := c.ShouldBind(&conf); err != nil {
		log.Error(err)
		c.JSON(200, model.HookReply{
			Code: model.ParseParamFail,
			Msg:  model.ParseParamFailMsg,
		})
		return
	}
	// do something
	conf.RemoteIp = c.RemoteIP()
	log.Info("收到zlm上线事件,media_server_id:", conf.GeneralMediaServerId, " ip:", conf.RemoteIp, " port:", conf.HttpPort)
	go service.Media().Online(conf)
	replyAllowMsg(c)
}

func (m MediaHookController) OnServerKeepalive(c *gin.Context) {
	keepalive := model.ServerKeepalive{}
	if err := c.ShouldBindJSON(&keepalive); err != nil {
		log.Error(err)
		c.JSON(200, model.HookReply{
			Code: model.ParseParamFail,
			Msg:  model.ParseParamFailMsg,
		})
		return
	}
	// do something
	log.Debugf("收到 Zlm id: %s 心跳", keepalive.MediaServerId)
	if err := service.Media().Keepalive(keepalive.MediaServerId); err != nil {
		log.Error(err)
	}
	_ = cron.StopTask(keepalive.MediaServerId, cron.TaskKeepLive)
	if err := cron.StartTask(keepalive.MediaServerId, cron.TaskKeepLive, 20*time.Second, func() {
		log.Debugf("ZLM心跳超时, id: %s", keepalive.MediaServerId)
		if err := service.Media().Offline(keepalive.MediaServerId); err != nil {
			log.Error("ZLM心跳超时, 更新数据失败, err: ", err)
		}
	}); err != nil {
		log.Error(err)
	}
	replyAllowMsg(c)
}

func (m MediaHookController) OnPlay(c *gin.Context) {
	hookParam := model.OnPlayHookParam{}
	if err := c.ShouldBindJSON(&hookParam); err != nil {
		log.Error(err)
		c.JSON(200, model.HookReply{
			Code: model.ParseParamFail,
			Msg:  model.ParseParamFailMsg,
		})
		return
	}
	// do something

	replyAllowMsg(c)
}

func (m MediaHookController) OnPublish(c *gin.Context) {
	hookParam := model.OnPublishHookParam{}
	if err := c.ShouldBindJSON(&hookParam); err != nil {
		log.Error(err)
		c.JSON(200, model.HookReply{
			Code: model.ParseParamFail,
			Msg:  model.ParseParamFailMsg,
		})
		return
	}
	// do something

	c.JSON(200, model.NewOnPublishDefaultReply())
}

func (m MediaHookController) OnStreamChanged(c *gin.Context) {
	hookParam := model.OnStreamChangedParam{}
	if err := c.ShouldBindJSON(&hookParam); err != nil {
		log.Error(err)
		c.JSON(200, model.HookReply{
			Code: model.ParseParamFail,
			Msg:  model.ParseParamFailMsg,
		})
		return
	}
	// do something
	log.Info("收到流改变事件,stream_id:", hookParam.Stream, " register: ", hookParam.Register, " protocol:", hookParam.Schema)
	replyAllowMsg(c)
}

func (m MediaHookController) OnStreamNOneReader(c *gin.Context) {
	closeStream := func() {
		c.JSON(200, model.OnStreamNoneReaderReply{
			Code:  0,
			Close: true,
		})
	}

	hookParam := model.OnStreamNoneReader{}
	if err := c.ShouldBindJSON(&hookParam); err != nil {
		log.Error(err)
		c.JSON(200, model.HookReply{
			Code: model.ParseParamFail,
			Msg:  model.ParseParamFailMsg,
		})
		return
	}

	// do something
	log.Info("收到流无人观看事件,stream_id:", hookParam.Stream, "media_server_id:", hookParam.MediaServerId)

	s := strings.Split(hookParam.Stream, "_")
	if len(s) == 0 {
		log.Error("stream split by '_' fail")
		closeStream()
		return
	}

	device, exist := service.Device().GetByDeviceId(s[0])
	if !exist {
		log.Error("not found device by deviceId query")
		closeStream()
		return
	}

	err := gbsip.StopPlay(hookParam.Stream, s[1], device)
	if err != nil {
		log.Errorf("%+v", err)
	}
	closeStream()
}

func (m MediaHookController) OnStreamNotFound(c *gin.Context) {

}

func (m MediaHookController) OnRtpServerTimeout(c *gin.Context) {

}

func (m MediaHookController) OnFlowReport(c *gin.Context) {

}

func (m MediaHookController) OnHttpAccess(c *gin.Context) {

}

func (m MediaHookController) OnRecordMp4(c *gin.Context) {

}

func (m MediaHookController) OnRtspAuth(c *gin.Context) {

}

func (m MediaHookController) OnRtspRealm(c *gin.Context) {

}

func (m MediaHookController) OnShellLogin(c *gin.Context) {

}

func replyAllowMsg(c *gin.Context) {
	c.JSON(200, model.HookReply{
		Code: model.RespondSuccess,
		Msg:  model.SuccessMsg,
	})
}
