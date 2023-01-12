package server

import (
	"github.com/chenjianhao66/go-GB28181/internal/controller"
	"github.com/chenjianhao66/go-GB28181/internal/store/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type apiServer struct {
	engine *gin.Engine
}

func newApiServer() *apiServer {
	return &apiServer{
		engine: gin.New(),
	}
}

func (a *apiServer) initRoute() {
	installController(a.engine)
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
