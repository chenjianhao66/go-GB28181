package controller

import (
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/gin-gonic/gin"
)

type MediaHookController struct{}

func NewMediaHookController() MediaHookController {
	return MediaHookController{}
}

// OnServerStarted 服务器启动事件，可以用于监听服务器崩溃重启；此事件对回复不敏感。
func (m MediaHookController) OnServerStarted(c *gin.Context) {
	defer c.JSON(200, "success")
	conf := &model.MediaConfig{}

	if err := c.ShouldBind(conf); err != nil {
		log.Error(err)
		c.JSON(200, model.HookReply{
			Code: model.ParseParamFail,
			Msg:  model.ParseParamFailMsg,
		})
		return
	}
	// do something
	conf.RemoteIp = c.RemoteIP()
	c.JSON(200, "")
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

	replyAllowMsg(c)
}

func (m MediaHookController) OnStreamNOneReader(c *gin.Context) {
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
	// TODO 发送bye消息给流发送者，关闭推流
	c.JSON(200, model.OnStreamNoneReaderReply{
		Code:  0,
		Close: true,
	})
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
		Code: model.Success,
		Msg:  model.SuccessMsg,
	})
}
