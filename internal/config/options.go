package config

import "github.com/spf13/viper"

type Options struct {
	ServerOptions *ServerOptions `json:"server" mapstructure:"server"`
	MySQLOptions  *MySQLOptions  `json:"mysql" mapstructure:"mysql"`
	SIPOptions    *SIPOptions    `json:"sip" mapstructure:"sip"`
	LogOptions    *LogOptions    `json:"log" mapstructure:"log"`
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
		SIPOptions:    newSIPOptions(),
		LogOptions:    newLogOptions(),
	}
}

func loadConfig() {
	viper.AddConfigPath("config/")
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("load config fail,please check your config file whether in config/ in the directory")
	}
}
