package cron

import (
	"github.com/chenjianhao66/go-GB28181/internal/log"
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
	select {
	case <-k.timer.C:
		log.Warnf("Device offline! DeviceId: %+v, at: %+v", k.deviceId, time.Now().GoString())
		k.runFunc()
		delete(taskList[k.deviceId], TaskKeepLive)
		return
	}
}
