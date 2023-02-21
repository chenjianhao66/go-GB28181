package controller

import (
	"github.com/chenjianhao66/go-GB28181/internal/gb"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/chenjianhao66/go-GB28181/internal/service"
	"github.com/gin-gonic/gin"
	"strings"
)

type MediaHookController struct{}

func NewMediaHookController() MediaHookController {
	return MediaHookController{}
}

// OnServerStarted 服务器启动事件，可以用于监听服务器崩溃重启；此事件对回复不敏感。
func (m MediaHookController) OnServerStarted(c *gin.Context) {
	log.Info("zlm 上线...")
	defer c.JSON(200, "success")
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
	log.Info("收到zlm上线事件,media_server_id:", conf.GeneralMediaServerId, "ip:", conf.RemoteIp, "port:", conf.HttpPort)
	conf.RemoteIp = c.RemoteIP()
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
	log.Info("收到流改变事件,stream_id:", hookParam.Stream, "register: ", hookParam.Register, "protocol:", hookParam.Schema)
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

	err := gb.SipCommand.StopPlay(hookParam.Stream, s[1], device)
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
