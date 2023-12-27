package gbserver

import (
	"context"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/controller"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/service"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage/sqlite"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

type apiServer struct {
	h      *http.Server
	engine *gin.Engine
	c      *apiConfig
}

type apiConfig struct {
	mediaOption  *option.MediaOptions
	serverOption *option.ServerOptions
	mysqlOption  *option.MySQLOptions
}

func newApiServer(config *apiConfig) *apiServer {
	return &apiServer{
		engine: gin.New(),
		c:      config,
	}
}

func (a *apiServer) initRoute() {
	a.installController()
	a.h = &http.Server{
		Handler: a.engine,
		Addr:    fmt.Sprintf(":%s", a.c.serverOption.Port),
	}
}

func (a *apiServer) Close() error {
	ctx := context.Background()
	withTimeout, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	log.Info("apiserver shutdown...")
	if err := a.h.Shutdown(withTimeout); err != nil {
		log.Info("close apiserver fail")
		panic(err)
	}
	return nil
}

func (a *apiServer) installController() {
	//store := mysql.GetMySQLFactory()
	// 使用sqlite3
	store := sqlite.GetSqliteFactory()
	service.InitService(store)
	initMediaHookRoute(a.engine.Group("/index/hook"))
	initDeviceRoute(a.engine.Group("/device"), store)
	initChannelRoute(a.engine.Group("/channel"), store)
	initControlRoute(a.engine.Group("/control"))
	initPlayRoute(a.engine.Group("/play"), store)
	initSwaggerRoute(a.engine.Group("/"))
}

func initSwaggerRoute(group *gin.RouterGroup) {
	group.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func initPlayRoute(group *gin.RouterGroup, store storage.Factory) {
	playController := controller.NewPlayController(store)
	group.POST("/start/:deviceId/:channelId", playController.Play)
}

func initControlRoute(group *gin.RouterGroup) {
	c := controller.NewControlController()
	group.POST("ptz", c.ControlPTZ)
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

	// 设备的基本配置
	group.POST("/config/basic", d.BasicParamsConfig)
	group.GET("/config/basic/:deviceId", d.BasicParamsQuery)
	// 查询设备状态
	group.GET("/status/:deviceId", d.StatusQuery)
	// 查询设备文件目录
	//group.GET("/catalog",)

	// 订阅
	group.POST("/subscribe/alarm:deviceId", d.AlarmSubscribe)
	group.POST("/subscribe/catalog:deviceId", d.CatalogSubscribe)
	group.POST("/subscribe/mobilePosition:deviceId", d.MobilePositionSubscribe)

}

func initChannelRoute(group *gin.RouterGroup, factory storage.Factory) {
	c := controller.NewChannelController(factory)
	group.GET("/list/:device", c.List)
}
