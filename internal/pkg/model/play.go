package model

import (
	"fmt"
	"strings"
)

// StreamInfo 流信息
type StreamInfo struct {
	// 承载该流信息的流媒体id
	MediaServerId string `json:"mediaServerId"`

	// 应用名
	App string `json:"app"`

	// ip地址
	Ip string `json:"ip"`

	// 输出流的设备id
	DeviceID string `json:"deviceID"`

	// 输出流的设备通道id
	ChannelId string `json:"channelId"`

	// 流名称
	Stream string `json:"stream"`

	// rtmp地址
	Rtmp string `json:"rtmp"`

	// rtsp地址
	Rtsp string `json:"rtsp"`

	// flv地址
	Flv string `json:"flv"`

	// https-flv地址
	HttpsFlv string `json:"httpsFlv"`
	//WSFly         string `json:"WSFly"`
	//WSSFly        string `json:"WSSFly"`

	// fmp4地址
	Fmp4 string `json:"fmp4"`
	//HttpsFmp4     string `json:"httpsFmp4"`
	//WSFmpt4       string `json:"WSFmpt4"`

	// hls地址
	Hls string `json:"hls"`
	//HttpsHls      string `json:"httpsHls"`
	//WsHls         string `json:"wsHls"`

	// ts地址
	Ts string `json:"ts"`
	//HttpsTs       string `json:"httpsTs"`
	//WebsocketTs   string `json:"websocketTs"`

	// 该流的ssrc
	Ssrc string `json:"ssrc"`
}

const (
	rtsp = "rtsp://%s:554/rtp/%s"
	rtmp = "rtmp://%s:1935/rtp/%s"
	http = "http://%s:8000/rtp/%s/hls.m3u8"
	flv  = "http://%s8000/rtp/%s.live.flv"
	fmp4 = "http://%s8000/rtp/%s.llive.mp4"
	ts   = "http://%s8000/rtp/%s.llive.ts"
)

func MustNewStreamInfo(mediaId, mediaIp, stream, ssrc string) StreamInfo {
	index := strings.Index(stream, "_")
	deviceId := stream[0:index]
	channelID := stream[index:]
	return StreamInfo{
		MediaServerId: mediaId,
		App:           "rtp",
		Ip:            mediaIp,
		DeviceID:      deviceId,
		ChannelId:     channelID,
		Stream:        stream,
		Rtmp:          fmt.Sprintf(rtmp, mediaIp, stream),
		Rtsp:          fmt.Sprintf(rtsp, mediaIp, stream),
		Hls:           fmt.Sprintf(http, mediaIp, stream),
		Flv:           fmt.Sprintf(flv, mediaIp, stream),
		Fmp4:          fmt.Sprintf(fmp4, mediaIp, stream),
		Ts:            fmt.Sprintf(ts, mediaIp, stream),
		Ssrc:          ssrc,
	}
}

// SipTransaction transaction info of sip request
type SipTransaction struct {
}
