package option

import (
	"github.com/spf13/pflag"
)

type SIPOptions struct {
	Ip        string `json:"ip,omitempty" mapstructure:"ip"`
	Port      string `json:"port,omitempty" mapstructure:"port"`
	Domain    string `json:"domain,omitempty" mapstructure:"domain"`
	Id        string `json:"id" mapstructure:"id"`
	Password  string `json:"password,omitempty" mapstructure:"password"`
	UserAgent string `json:"user-agent" mapstructure:"user-agent"`
}

func NewSIPOptions() *SIPOptions {
	return &SIPOptions{
		Ip:     "127.0.0.1",
		Port:   "5060",
		Domain: "4401020049",
		Id:     "44010200492000000001",
	}

}

func (s *SIPOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&s.Ip, "sip.ip", s.Ip, "sip服务器的地址，填本机ip")
	fss.StringVar(&s.Port, "sip.port", s.Port, "sip服务监听的端口")
	fss.StringVar(&s.Domain, "sip.domain", s.Domain, "sip服务器国标域编码")
	fss.StringVar(&s.Id, "sip.id", s.Id, "sip服务器国标唯一编码")
}
