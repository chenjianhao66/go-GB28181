package util

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model/constant"
	"sync"
)

type ssrc struct {
	m         sync.Mutex
	isUsed    []string
	isNotUsed []string
}

func (s *ssrc) getSSRC(t SsrcPrefix) string {
	s.m.Lock()
	defer s.m.Unlock()
	ssrc := s.isNotUsed[0]
	s.isUsed = append(s.isUsed, ssrc)
	s.isNotUsed = s.isNotUsed[1:]
	key := fmt.Sprintf("%d%s%s", t, config.SIPDomain()[3:8], ssrc)
	return key
}

type SsrcPrefix int

const (
	RealTime SsrcPrefix = iota
	History
)

var ssrcInfo = initSSRC()

func initSSRC() *ssrc {
	var noUsed []string
	for i := 1; i < constant.MaxStreamCount; i++ {
		ssrc := fmt.Sprintf("%04d", i)
		noUsed = append(noUsed, ssrc)
	}
	return &ssrc{
		m:         sync.Mutex{},
		isUsed:    make([]string, 0, constant.MaxStreamCount),
		isNotUsed: noUsed,
	}
}

func GetSSRC(t SsrcPrefix) string {
	return ssrcInfo.getSSRC(t)
}
