package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Options struct {
	ServerOptions *ServerOptions `json:"server" mapstructure:"server"`
	MySQLOptions  *MySQLOptions  `json:"mysql" mapstructure:"mysql"`
	RedisOptions  *RedisOptions  `json:"redis" mapstructure:"redis"`
	SIPOptions    *SIPOptions    `json:"sip" mapstructure:"sip"`
	LogOptions    *LogOptions    `json:"log" mapstructure:"log"`
	MediaOptions  *MediaOptions  `json:"media" mapstructure:"media"`
}

const defaultConfigName = "application"

var (
	o = initOptions()
)

func initOptions() *Options {
	loadConfig()
	return &Options{
		ServerOptions: NewServerOptions(),
		MySQLOptions:  NewMySQLOptions(),
		RedisOptions:  newRedisOptions(),
		SIPOptions:    newSIPOptions(),
		LogOptions:    newLogOptions(),
		MediaOptions:  newMediaOption(),
	}
}

func loadConfig() {
	viper.AddConfigPath("config/")
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType("yml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Print("load config fail,please check your config file whether in config/ in the directory")
	}
}
