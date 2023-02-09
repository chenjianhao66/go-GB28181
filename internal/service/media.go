package service

import (
	"encoding/json"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/chenjianhao66/go-GB28181/internal/model/constant"
	"github.com/chenjianhao66/go-GB28181/internal/storage"
	"github.com/chenjianhao66/go-GB28181/internal/storage/cache"
)

type IMedia interface {
	Online(config model.MediaConfig)
}

type mediaService struct {
	store storage.Factory
}

var mService = new(mediaService)

func Media() IMedia {
	return mService
}

func (m *mediaService) Online(c model.MediaConfig) {
	newMediaDetail := model.NewMediaDetailWithConfig(&c)
	if err := m.store.Media().Save(newMediaDetail); err != nil {
		log.Error(err)
		return
	}
	// please check this stream server if Whether in cache
	key := fmt.Sprintf("%s_%s", constant.MediaServerPrefix, newMediaDetail.ID)
	cacheDetail, _ := cache.Get(key)
	if cacheDetail == "" {
		newMediaDetail.SsrcConfig = model.NewSsrcConfig(newMediaDetail.ID, config.SIPDomain())
	} else {
		oldMediaDetail := model.MediaDetail{}
		err := json.Unmarshal([]byte(cacheDetail.(string)), &oldMediaDetail)
		if err != nil {
			log.Error("JSON数据解析到结构体失败!", err)
			return
		}
		newMediaDetail.SsrcConfig = oldMediaDetail.SsrcConfig
	}
	cache.Set(key, newMediaDetail)
	log.Info(fmt.Sprintf("ZleMedia流媒体连接成功,id: [ %s ] , addr: [ %s:%v ]", newMediaDetail.ID, newMediaDetail.Ip, newMediaDetail.HttpPort))
}
