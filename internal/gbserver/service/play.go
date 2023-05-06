package service

import (
	"encoding/json"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage/cache"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model/constant"
	"github.com/pkg/errors"
)

type IPlay interface {
	Play(deviceId, channelId string) (model.StreamInfo, error)
}

type playService struct{}

var (
	deviceNotFound = errors.New("device not found")
)

func Play() IPlay {
	return playService{}
}

func (p playService) Play(deviceId, channelId string) (model.StreamInfo, error) {
	var (
		streamInfo  model.StreamInfo
		mediaDetail model.MediaDetail
		streamId    = fmt.Sprintf("%s_%s", deviceId, channelId)
	)
	device, ok := Device().GetByDeviceId(deviceId)
	if !ok {
		return model.StreamInfo{}, deviceNotFound
	}

	mediaDetail, err := Media().GetDefaultMedia()
	if err != nil {
		return model.StreamInfo{}, err
	}

	key := fmt.Sprintf("%s:%s", constant.StreamInfoPrefix, streamId)
	streamJSON, _ := cache.Get(key)

	if streamJSON != "" {
		if err := json.Unmarshal([]byte(streamJSON.(string)), &streamInfo); err != nil {
			return model.StreamInfo{}, errors.WithMessage(err, "unmarshal json data to struct fail")
		}

		rtpServerInfo, err := Media().GetRtpServerInfo(streamId, mediaDetail)
		if err != nil {
			return model.StreamInfo{}, err
		}

		if rtpServerInfo.Code == model.RespondSuccess {
			if rtpServerInfo.Exist == true {
				return streamInfo, nil
			} else {
				streamInfo = model.StreamInfo{}
			}
		} else {
			// zlm服务连接失败，重新创建rtp服务并连接
			log.Errorf("media api response: %+v\n", rtpServerInfo.Msg)
			streamInfo = model.StreamInfo{}
		}
	}

	// 判断流信息对象是否是默认值，是默认值的话代表没有这个流信息或者rtp服务连接失败
	if streamInfo == (model.StreamInfo{}) {
		rtpPort, ssrc, err := Media().OpenRtpServer(mediaDetail, streamId)
		streamInfo.Ssrc = ssrc

		if err != nil {
			return model.StreamInfo{}, errors.WithMessage(err, "create rtp server fail")
		}
		return gbsip.Play(device, mediaDetail, streamId, ssrc, channelId, rtpPort)
	}

	return model.StreamInfo{}, nil
}
