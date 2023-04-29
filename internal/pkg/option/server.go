package option

import (
	"github.com/spf13/pflag"
)

type ServerOptions struct {
	Port string `json:"port" mapstructure:"port"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		Port: "18080",
	}
}

func (s *ServerOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&s.Port, "server.port", s.Port, "gb服务器的http端口")
}
