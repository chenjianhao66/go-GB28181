package model

import (
	"fmt"
	"strings"
)

// StreamInfo 流信息
type StreamInfo struct {
	MediaServerId string `json:"mediaServerId"`
	App           string `json:"app"`
	Ip            string `json:"ip"`
	DeviceID      string `json:"deviceID"`
	ChannelId     string `json:"channelId"`
	Stream        string `json:"stream"`
	Rtmp          string `json:"rtmp"`
	Rtsp          string `json:"rtsp"`
	Flv           string `json:"flv"`
	HttpsFlv      string `json:"httpsFlv"`
	//WSFly         string `json:"WSFly"`
	//WSSFly        string `json:"WSSFly"`
	Fmp4 string `json:"fmp4"`
	//HttpsFmp4     string `json:"httpsFmp4"`
	//WSFmpt4       string `json:"WSFmpt4"`
	Hls string `json:"hls"`
	//HttpsHls      string `json:"httpsHls"`
	//WsHls         string `json:"wsHls"`
	Ts string `json:"ts"`
	//HttpsTs       string `json:"httpsTs"`
	//WebsocketTs   string `json:"websocketTs"`
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
