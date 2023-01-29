package server

import (
	"context"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/controller"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/store/mysql"
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
	withTimeout, _ := context.WithTimeout(ctx, 1*time.Second)
	if err := a.h.Shutdown(withTimeout); err != nil {
		log.Info("close apiserver fail")
		panic(err)
	}
	log.Info("apiserver shutdown...")
	return nil
}

func installController(g *gin.Engine) *gin.Engine {
	store, _ := mysql.GetMySQLFactory()
	g.GET("version", func(context *gin.Context) {
		m := map[string]string{
			"data": "o",
		}
		context.JSON(http.StatusOK, m)
	})

	device := g.Group("/device")
	{
		// TODO 后续添加接口支持
		deviceController := controller.NewDeviceController(store)
		device.GET("list", deviceController.List)
	}

	return g
}
