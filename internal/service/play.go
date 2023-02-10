package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/model"
	"github.com/chenjianhao66/go-GB28181/internal/model/constant"
	"github.com/chenjianhao66/go-GB28181/internal/storage/cache"
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
		streamInfo model.StreamInfo
	)
	device, ok := Device().GetByDeviceId(deviceId)
	log.Info(device)

	if !ok {
		return model.StreamInfo{}, deviceNotFound
	}
	// TODO 从缓存中获取流信息是否存在，存在则直接返回
	streamId := fmt.Sprintf("%s:%s", constant.StreamInfoPrefix, deviceId+"_"+channelId)
	streamJSON, _ := cache.Get(streamId)

	if streamJSON != "" {
		if err := json.Unmarshal([]byte(streamJSON.(string)), &streamInfo); err != nil {
			log.Error("解析到流信息结构体失败")
			return model.StreamInfo{}, errors.New("点播失败")
		}
		mediaDetail, _ := Media().GetMedia(streamInfo.MediaServerId)

		Media().GetRtpServerInfo(streamId, mediaDetail)

	}

	//if err != nil {
	//
	//}
	// TODO
	return model.StreamInfo{}, nil
}

func getRtpInfo() {

}
