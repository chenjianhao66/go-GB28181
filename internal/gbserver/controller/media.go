package controller

import (
	srv "github.com/chenjianhao66/go-GB28181/internal/gbserver/service"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
	"github.com/gin-gonic/gin"
)

type MediaController struct {
	srv srv.Service
}

func NewMediaController(store storage.Factory) *MediaController {
	return &MediaController{
		srv: srv.NewService(store),
	}
}

func (m *MediaController) List(ctx *gin.Context) {
	list, err := m.srv.Media().List()

	if err != nil {
		newResponse(ctx).fail("获取流媒体列表失败")
		return
	}
	newResponse(ctx).successWithAny(list)
}
