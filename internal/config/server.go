package config

import "github.com/spf13/viper"

type ServerOptions struct {
	Port string `json:"port" mapstructure:"port"`
}

func NewServerOptions() *ServerOptions {
	s := new(ServerOptions)
	_ = viper.UnmarshalKey("server", s)
	return s
}
