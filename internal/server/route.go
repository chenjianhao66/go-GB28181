package server

import (
	"context"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/controller"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/storage"
	"github.com/chenjianhao66/go-GB28181/internal/storage/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type apiServer struct {
	h      *http.Server
	engine *gin.Engine
}

func newApiServer() *apiServer {
	return &apiServer{
		engine: gin.New(),
	}
}

func (a *apiServer) initRoute() {
	installController(a.engine)
	a.h = &http.Server{
		Handler: a.engine,
		Addr:    fmt.Sprintf(":%s", config.ServerPort()),
	}
}

func (a *apiServer) Close() error {
	ctx := context.Background()
	withTimeout, cancelFunc := context.WithTimeout(ctx, 1*time.Second)
	defer cancelFunc()
	if err := a.h.Shutdown(withTimeout); err != nil {
		log.Info("close apiserver fail")
		panic(err)
	}
	log.Info("apiserver shutdown...")
	return nil
}

func installController(g *gin.Engine) *gin.Engine {
	store := mysql.GetMySQLFactory()
	initMediaHookRoute(g.Group("/index/hook"))
	initDeviceRoute(g.Group("/device"), store)
	initChannelRoute(g.Group("/channel"), store)
	return g
}

func initMediaHookRoute(group *gin.RouterGroup) {
	hook := controller.NewMediaHookController()
	group.POST("on_server_started", hook.OnServerStarted)
	group.POST("on_server_keepalive", hook.OnServerKeepalive)
	group.POST("on_play", hook.OnPlay)
	group.POST("on_publish", hook.OnPublish)
	group.POST("on_stream_changed", hook.OnStreamChanged)
	group.POST("on_stream_none_reader", hook.OnStreamNOneReader)
	group.POST("on_stream_not_found", hook.OnStreamNotFound)
	group.POST("on_rtp_server_timeout", hook.OnRtpServerTimeout)
	group.POST("on_flow_report", hook.OnFlowReport)
	group.POST("on_http_access", hook.OnHttpAccess)
	group.POST("on_record_mp4", hook.OnRecordMp4)
	group.POST("on_rtsp_auth", hook.OnRtspAuth)
	group.POST("on_rtsp_realm", hook.OnRtspRealm)
	group.POST("on_shell_login", hook.OnShellLogin)
}

func initDeviceRoute(group *gin.RouterGroup, factory storage.Factory) {
	d := controller.NewDeviceController(factory)
	group.GET("/list", d.List)
}

func initChannelRoute(group *gin.RouterGroup, factory storage.Factory) {
	c := controller.NewChannelController(factory)
	group.GET("/list/:device", c.List)
}
