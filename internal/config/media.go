package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type MediaOptions struct {
	Id       string `json:"id,omitempty" mapstructure:"id"`
	Ip       string `json:"ip,omitempty" mapstructure:"ip"`
	HttpPort string `json:"http-port,omitempty" mapstructure:"http-port"`
	Secret   string `json:"secret,omitempty" mapstructure:"secret"`
}

func newMediaOption() *MediaOptions {
	m := new(MediaOptions)
	err := viper.UnmarshalKey("media", m)
	fmt.Printf("%+v", m)
	if err != nil {
		panic("parser config 'media item fail:" + err.Error())
	}
	return m
}

func MediaSecret() string {
	return o.MediaOptions.Secret
}

func MediaIp() string {
	return o.MediaOptions.Ip
}

func MediaServiceId() string {
	return o.MediaOptions.Id
}

func MediaHttpPort() string {
	return o.MediaOptions.HttpPort
}
