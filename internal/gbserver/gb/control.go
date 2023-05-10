package gb

import (
	"encoding/xml"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/parser"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/syn"
	"github.com/ghettovoice/gosip/sip"
)

func deviceConfigQueryHandler(req sip.Request, tx sip.ServerTransaction) {
	log.Debugf("获取到的configDownload消息：\n%s", req.Body())
	defer func() {
		_ = responseAck(tx, req)
	}()

	cfg := &gbsip.DeviceBasicConfigResp{}

	if err := xml.Unmarshal([]byte(req.Body()), cfg); err != nil {
		b, err := gbkToUtf8([]byte(req.Body()))
		if err != nil {
			log.Error(err)
			return
		}
		err = xml.Unmarshal(b, cfg)
		if err != nil {
			log.Error(err)
			return
		}
	}

	if cfg.R.Result != "OK" {
		return
	}

	syn.HasSyncTask(fmt.Sprintf("%s_%s", syn.KeyControlDeviceConfigQuery, cfg.DeviceID.DeviceID), func(e *syn.Entity) {
		e.Ok(cfg)
	})

	_ = storage.updateDeviceBasicConfig(*cfg)
}

func deviceConfigResponseHandler(req sip.Request, tx sip.ServerTransaction) {
	r := parser.GetResultFromXML(req.Body())
	if r == "" {
		log.Error("获取不到响应信息中的Result字段")
		return
	}

	if r == "ERROR" {
		log.Error("发送修改配置请求失败，请检查")
	} else {
		log.Debug("发送修改配置请求成功")
	}
}
