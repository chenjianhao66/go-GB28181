package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage/cache"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model/constant"
	util2 "github.com/chenjianhao66/go-GB28181/internal/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type IMedia interface {
	Online(config model.MediaConfig)
	Offline(mediaId string) error
	GetRtpServerInfo(stream string, mediaDetail model.MediaDetail) (model.GetRtpInfoResp, error)
	OpenRtpServer(detail model.MediaDetail, stream string) (rtpPort int, ssrc string, err error)
	GetMedia(serverId string) (model.MediaDetail, error)
	GetDefaultMedia() (model.MediaDetail, error)
	List() ([]model.MediaDetail, error)
	Keepalive(mediaId string) error
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
	if cast.ToString(viper.Get("media.ip")) == newMediaDetail.Ip {
		newMediaDetail.Default = true
	}
	if err := m.store.Media().Save(newMediaDetail); err != nil {
		log.Error(err)
		return
	}
	// please check this stream server if Whether in cache
	key := fmt.Sprintf("%s:%s", constant.MediaServerPrefix, newMediaDetail.ID)
	cacheDetail, _ := cache.Get(key)
	if cacheDetail == nil {
		//newMediaDetail.SsrcConfig = model.NewSsrcConfig(newMediaDetail.ID, config.SIPDomain())
	} else {
		oldMediaDetail := model.MediaDetail{}
		if err := json.Unmarshal(cacheDetail.([]byte), &oldMediaDetail); err != nil {
			log.Errorf("unmarshal json data to struct fail, err: %v", err)
		}
		newMediaDetail.SsrcConfig = oldMediaDetail.SsrcConfig
	}
	cache.Set(key, newMediaDetail)
	log.Info(fmt.Sprintf("ZleMedia流媒体连接成功,id: [ %s ] , addr: [ %s:%v ]", newMediaDetail.ID, newMediaDetail.Ip, newMediaDetail.HttpPort))
}

func (m *mediaService) Offline(mediaId string) error {
	media, err := m.GetMedia(mediaId)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("根据id [ %s ] 获取流媒体信息失败", mediaId))
	}
	media.Status = false
	return m.store.Media().Save(media)
}

// GetRtpServerInfo 从流媒体服务获取rtp明细信息
func (m *mediaService) GetRtpServerInfo(stream string, mediaDetail model.MediaDetail) (model.GetRtpInfoResp, error) {
	params := map[string]interface{}{
		"secret":    mediaDetail.Secret,
		"stream_id": stream,
	}

	url := fmt.Sprintf(constant.MediaGetRtpInfoApiUrl, mediaDetail.Ip, mediaDetail.HttpPort)

	result, err := util2.SendPost(url, params)
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
	ssrc = util2.GetSSRC(util2.RealTime)

	// todo 先写死ip为localhost
	url := fmt.Sprintf(constant.MediaCreateRtpApiUrl, "localhost", detail.HttpPort)
	params := map[string]interface{}{
		"secret":     detail.Secret,
		"port":       0,
		"enable_tcp": 1,
		"stream_id":  stream,
	}
	body, err := util2.SendPost(url, params)
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

func (m *mediaService) GetDefaultMedia() (model.MediaDetail, error) {
	key := fmt.Sprintf("%s:%s", constant.MediaServerPrefix, cast.ToString(viper.Get("media.id")))
	data, err := cache.Get(key)
	if err != nil {
		return model.MediaDetail{}, errors.WithMessage(err, "get media detail from cache fail")
	}

	if data != "" {
		detail := model.MediaDetail{}
		err := json.Unmarshal(data.([]byte), &detail)
		if err != nil {
			return model.MediaDetail{}, errors.WithMessage(err, "unmarshal json data to struct fail")
		}
		return detail, nil
	}

	// 从数据库获取
	detail, err := m.store.Media().GetMediaByID(cast.ToString(viper.Get("media.id")))
	if err != nil {
		return model.MediaDetail{}, errors.WithMessage(err, "get media detail from database fail")
	}

	return detail, nil
}

func (m *mediaService) List() ([]model.MediaDetail, error) {
	return m.store.Media().List()
}

func (m *mediaService) Keepalive(mediaId string) error {
	detail, err := m.store.Media().GetMediaByID(mediaId)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("根据id [ %s ] 获取流媒体信息失败", mediaId))
	}
	detail.LastKeepaliveTime = time.Now()
	detail.Status = true
	return m.store.Media().Save(detail)
}
