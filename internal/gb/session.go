package gb

import (
	"encoding/json"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/model/constant"
	"github.com/chenjianhao66/go-GB28181/internal/storage/cache"
	"github.com/pkg/errors"
)

// SipTX sip会话事务
type SipTX struct {
	DeviceId  string `json:"deviceId,omitempty"`
	ChannelId string `json:"channelId,omitempty"`
	SSRC      string `json:"SSRC,omitempty"`
	CallId    string `json:"callId,omitempty"`
	FromTag   string `json:"fromTag,omitempty"`
	ToTag     string `json:"toTag,omitempty"`
	ViaBranch string `json:"viaBranch,omitempty"`
}

type txManage struct{}

var streamSessionManage txManage

// 保存sip事务信息
func (s txManage) saveStreamSession(deviceId string, channelId string, ssrc string, callId string, fromTag string, toTag string, viaBranch string) {
	tx := SipTX{
		DeviceId:  deviceId,
		ChannelId: channelId,
		SSRC:      ssrc,
		CallId:    callId,
		FromTag:   fromTag,
		ToTag:     toTag,
		ViaBranch: viaBranch,
	}

	key := fmt.Sprintf("%s:%s", constant.StreamTransactionPrefix, deviceId+"_"+channelId)
	cache.Set(key, tx)
}

func (s txManage) getTx(deviceId, channelId string) (SipTX, error) {
	key := fmt.Sprintf("%s:%s", constant.StreamTransactionPrefix, deviceId+"_"+channelId)
	j, err := cache.Get(key)

	if err != nil {
		return SipTX{}, errors.WithMessage(err, "gei sip session tx fail")
	}

	var tx SipTX

	err = json.Unmarshal([]byte(j.(string)), &tx)
	if err != nil {
		return SipTX{}, errors.WithMessage(err, "unmarshal json data to struct fail")
	}

	return tx, nil
}
