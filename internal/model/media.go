package model

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/model/constant"
	"github.com/spf13/cast"
	"time"
)

// MediaDetail 流媒体明细
type MediaDetail struct {
	// secret
	ID                string     `gorm:"column:id;primaryKey;unique"`
	Ip                string     `gorm:"column:ip;size:100"`
	HookIp            string     `gorm:"column:hookIp"`
	SdpIp             string     `gorm:"column:sdpIp"`
	StreamIp          string     `gorm:"column:streamIp"`
	HttpPort          int        `gorm:"column:httpPort"`
	HttpSSlPort       int        `gorm:"column:httpSSLPort"`
	RtmpPort          int        `gorm:"column:rtmpPort"`
	RtmpSSlPort       int        `gorm:"column:rtmpSSLPort"`
	RtpProxyPort      int        `gorm:"column:rtpProxyPort"`
	RtspPort          int        `gorm:"column:rtspPort"`
	RtspSSLPort       int        `gorm:"column:rtspSSLPort"`
	RtpEnable         bool       `gorm:"column:rtpEnable"`
	RtpPortRange      string     `gorm:"column:rtpPortRange"`
	Secret            string     `gorm:"column:secret"`
	Default           bool       `gorm:"column:default"`
	HookAliveInterval int        `gorm:"column:hookAliveInterval"`
	CreateTime        time.Time  `gorm:"column:createTime;autoUpdateTime:milli"`
	UpdateTime        time.Time  `gorm:"column:updateTime;autoUpdateTime:milli"`
	LastKeepaliveTime time.Time  `gorm:"column:lastKeepaliveTime;autoUpdateTime:milli"`
	SsrcConfig        SsrcConfig `gorm:"-"`
}

func (m MediaDetail) TableName() string {
	return "mediaDetail"
}

func NewMediaDetailWithConfig(config *MediaConfig) MediaDetail {
	return MediaDetail{
		ID:                config.GeneralMediaServerId,
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
