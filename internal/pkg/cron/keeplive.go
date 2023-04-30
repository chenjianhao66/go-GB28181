package cron

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"time"
)

type keepLiveTask struct {
	timer    *time.Ticker
	deviceId string
	duration time.Duration
	runFunc  runFunc
}

func (k *keepLiveTask) start() error {
	k.timer = time.NewTicker(k.duration)
	go k.watch()

	return nil
}

func (k *keepLiveTask) cancel() error {
	k.timer.Stop()
	return nil
}

func (k *keepLiveTask) refresh() {
	k.timer.Reset(k.duration)
}

func (k *keepLiveTask) watch() {
	<-k.timer.C
	log.Warnf("设备离线！ 设备号: %s, 时间: %s", k.deviceId, time.Now().String())
	k.runFunc()
	taskList.deleteOneTask(k.deviceId, TaskKeepLive)
	return
}
