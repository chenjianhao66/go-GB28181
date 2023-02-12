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
	"github.com/chenjianhao66/go-GB28181/internal/util"
	"github.com/pkg/errors"
)

type IMedia interface {
	Online(config model.MediaConfig)
	GetRtpServerInfo(stream string, mediaDetail model.MediaDetail) (model.GetRtpInfoResp, error)
	OpenRtpServer(detail model.MediaDetail, stream string) (rtpPort int, ssrc string, err error)
	GetMedia(serverId string) (model.MediaDetail, error)
}

type mediaService struct {
	store storage.Factory
}

var mService = new(mediaService)

func Media() IMedia {
	return mService
}

// Online 流媒体服务上线事件
func (m *mediaService) Online(c model.MediaConfig) {
	newMediaDetail := model.NewMediaDetailWithConfig(&c)
	if err := m.store.Media().Save(newMediaDetail); err != nil {
		log.Error(err)
		return
	}
	// please check this stream server if Whether in cache
	key := fmt.Sprintf("%s:%s", constant.MediaServerPrefix, newMediaDetail.ID)
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

// GetRtpServerInfo 从流媒体服务获取rtp明细信息
func (m *mediaService) GetRtpServerInfo(stream string, mediaDetail model.MediaDetail) (model.GetRtpInfoResp, error) {
	params := map[string]interface{}{
		"secret":    config.MediaSecret(),
		"stream_id": stream,
	}

	url := fmt.Sprintf(constant.MediaGetRtpInfoApiUrl, mediaDetail.Ip, mediaDetail.HttpPort)

	result, err := util.SendPost(url, params)
	if err != nil {
		return model.GetRtpInfoResp{}, errors.WithMessage(err, "query media rtp server fail")
	}

	resp := model.GetRtpInfoResp{}
	err = json.Unmarshal([]byte(result), &resp)
	if err != nil {
		return model.GetRtpInfoResp{}, errors.WithMessage(err, "unmarshal data to struct fail")
	}

	return resp, err
}

// OpenRtpServer 创建rtp服务
func (m *mediaService) OpenRtpServer(detail model.MediaDetail, stream string) (rtpPort int, ssrc string, err error) {
	ssrc = util.GetSSRC(util.RealTime)

	url := fmt.Sprintf(constant.MediaCreateRtpApiUrl, detail.Ip, detail.HttpPort)
	params := map[string]interface{}{
		"secret":     config.MediaSecret(),
		"port":       0,
		"enable_tcp": 1,
		"stream_id":  stream,
	}
	body, err := util.SendPost(url, params)
	if err != nil {
		return 0, "", errors.WithMessage(err, "create rtp server fail")
	}

	resp := model.CreateRtpServerResp{}
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return 0, "", errors.WithMessage(err, "unmarshal data to struct fail")
	}

	if resp.Code != model.RespondSuccess {
		return 0, "", errors.New(resp.Msg)
	}

	rtpPort = resp.Port
	return
}

// GetMedia 从缓存里面获取一个流媒体明细
func (m *mediaService) GetMedia(serverId string) (model.MediaDetail, error) {
	j, err := cache.Get(serverId)
	if err != nil {
		return model.MediaDetail{}, errors.WithMessage(err, "GetMedia function happen error")
	}
	var detail = model.MediaDetail{}
	err = json.Unmarshal([]byte(j.(string)), &detail)
	if err != nil {
		return model.MediaDetail{}, errors.WithMessage(err, "GetMedia function unmarshal data to struct fail")
	}
	return detail, err
}
