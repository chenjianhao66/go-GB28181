package gbclient

import (
	"bytes"
	"context"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/gbsip"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/ghettovoice/gosip/sip"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	server struct {
		gb  *gbsip.Server
		ctx context.Context
	}
)

func newServer(c *ctlOption) *server {
	config := &gbsip.SipConfig{
		SipOption:   c.Sip,
		MysqlOption: nil,
		HandlerMap:  createHandlerMap(),
	}
	return &server{
		gb:  gbsip.NewServer(config),
		ctx: context.Background(),
	}
}

func (s *server) run() error {
	eg, ctx := errgroup.WithContext(s.ctx)
	defer ctx.Done()

	eg.Go(func() error {
		if err := s.gb.ListenUDP(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		return nil
	})

	eg.Go(func() error {
		if err := s.gb.ListenTCP(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		return nil
	})

	gbsip.Register(initSIPDeviceInfo())

	return eg.Wait()
}

func createHandlerMap() gbsip.RequestHandlerMap {
	m := make(map[sip.RequestMethod]func(req sip.Request, tx sip.ServerTransaction))
	m[sip.REGISTER] = registerHandler
	//m[sip.MESSAGE] = MessageHandler
	return m
}

func initSIPDeviceInfo() (devices []model.Device) {
	used := viper.ConfigFileUsed()
	//flag := pflag.Lookup("config")
	cfg, err := os.Open(used)
	if err != nil {
		log.Errorf("获取配置文件失败，err: %v", err)
		os.Exit(1)
	}
	bf := new(bytes.Buffer)
	_, err = bf.ReadFrom(cfg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	clientOpt := &ctlOption{}
	if err = yaml.Unmarshal(bf.Bytes(), clientOpt); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	//devices = make([]model.Device, 2)

	var cmap map[string]interface{}

	if err = yaml.Unmarshal(bf.Bytes(), &cmap); err != nil {
		os.Exit(1)
	}

	srv := model.Device{
		DeviceId: clientOpt.Sip.Id,
		Domain:   clientOpt.Sip.Domain,
		Ip:       clientOpt.Sip.Ip,
		Port:     clientOpt.Sip.Port,
	}
	devices = append(devices, srv)

	client := model.Device{
		DeviceId:  clientOpt.ClientOption.Id,
		Domain:    clientOpt.ClientOption.Domain,
		Ip:        clientOpt.ClientOption.Ip,
		Port:      clientOpt.ClientOption.Port,
		Transport: clientOpt.ClientOption.Transport,
	}
	devices = append(devices, client)
	return devices
}
