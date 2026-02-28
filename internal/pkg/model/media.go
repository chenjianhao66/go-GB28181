package model

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model/constant"
	"github.com/spf13/cast"
	"time"
)

// MediaDetail 流媒体明细
type MediaDetail struct {
	// secret
	ID                string     `gorm:"column:id;primaryKey;unique" json:"ID,omitempty"`
	Ip                string     `gorm:"column:ip;size:100" json:"ip,omitempty"`
	HookIp            string     `gorm:"column:hookIp" json:"hookIp,omitempty"`
	SdpIp             string     `gorm:"column:sdpIp" json:"sdpIp,omitempty"`
	StreamIp          string     `gorm:"column:streamIp" json:"streamIp,omitempty"`
	HttpPort          int        `gorm:"column:httpPort" json:"httpPort,omitempty"`
	HttpSSlPort       int        `gorm:"column:httpSSLPort" json:"httpSSlPort,omitempty"`
	RtmpPort          int        `gorm:"column:rtmpPort" json:"rtmpPort,omitempty"`
	RtmpSSlPort       int        `gorm:"column:rtmpSSLPort" json:"rtmpSSlPort,omitempty"`
	RtpProxyPort      int        `gorm:"column:rtpProxyPort" json:"rtpProxyPort,omitempty"`
	RtspPort          int        `gorm:"column:rtspPort" json:"rtspPort,omitempty"`
	RtspSSLPort       int        `gorm:"column:rtspSSLPort" json:"rtspSSLPort,omitempty"`
	RtpEnable         bool       `gorm:"column:rtpEnable" json:"rtpEnable,omitempty"`
	RtpPortRange      string     `gorm:"column:rtpPortRange" json:"rtpPortRange,omitempty"`
	Secret            string     `gorm:"column:secret" json:"secret,omitempty"`
	Default           bool       `gorm:"column:default" json:"default,omitempty"`
	Status            bool       `gorm:"column:status" json:"status"`
	HookAliveInterval int        `gorm:"column:hookAliveInterval" json:"hookAliveInterval,omitempty"`
	CreateTime        time.Time  `gorm:"column:createTime;autoUpdateTime:milli" json:"createTime"`
	UpdateTime        time.Time  `gorm:"column:updateTime;autoUpdateTime:milli" json:"updateTime"`
	LastKeepaliveTime time.Time  `gorm:"column:lastKeepaliveTime;autoUpdateTime:milli" json:"lastKeepaliveTime"`
	SsrcConfig        SsrcConfig `gorm:"-" json:"ssrcConfig"`
}

func (m MediaDetail) TableName() string {
	return "mediaDetail"
}

func NewMediaDetailWithConfig(config *MediaConfig) MediaDetail {
	return MediaDetail{
		ID:                config.GeneralMediaServerId,
		HookIp:            config.RemoteIp,
		Ip:                config.RemoteIp,
		SdpIp:             config.RemoteIp,
		StreamIp:          config.RemoteIp,
		HttpPort:          cast.ToInt(config.HttpPort),
		HttpSSlPort:       cast.ToInt(config.HttpSSLPort),
		RtmpPort:          cast.ToInt(config.RtmpPort),
		RtmpSSlPort:       cast.ToInt(config.RtmpSSLPort),
		RtpProxyPort:      cast.ToInt(config.RtpProxyPort),
		RtspPort:          cast.ToInt(config.RtspPort),
		RtspSSLPort:       cast.ToInt(config.RtspSSLPort),
		RtpEnable:         true,
		RtpPortRange:      config.RtpProxyPortRange,
		Secret:            config.ApiSecret,
		HookAliveInterval: cast.ToInt(config.HookAliveInterval),
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
		LastKeepaliveTime: time.Now(),
	}
}

// SsrcConfig ssrc配置
type SsrcConfig struct {
	MediaServerId string
	SsrcPrefix    string
	IsUsed        []string
	NotUsed       []string
}

// NewSsrcConfig 初始化一个ssrc配置
func NewSsrcConfig(mediaServerId, domain string) SsrcConfig {
	var noUsed []string
	for i := 1; i < constant.MaxStreamCount; i++ {
		var ssrc string
		if i < 10 {
			ssrc = fmt.Sprintf("000%d", i)
		} else if i < 100 {
			ssrc = fmt.Sprintf("00%d", i)
		} else if i < 1000 {
			ssrc = fmt.Sprintf("0%d", i)
		} else {
			ssrc = cast.ToString(i)
		}
		noUsed = append(noUsed, ssrc)
	}
	return SsrcConfig{
		MediaServerId: mediaServerId,
		SsrcPrefix:    domain[3:8],
		NotUsed:       noUsed,
		IsUsed:        make([]string, 0),
	}

}
