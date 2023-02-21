package config

import (
	"github.com/spf13/viper"
)

type SIPOptions struct {
	Ip        string `json:"ip,omitempty" mapstructure:"ip"`
	Port      string `json:"port,omitempty" mapstructure:"port"`
	Domain    string `json:"domain,omitempty" mapstructure:"domain"`
	Id        string `json:"id" mapstructure:"id"`
	Password  string `json:"password,omitempty" mapstructure:"password"`
	UserAgent string `json:"user-agent" mapstructure:"user-agent"`
}

func newSIPOptions() *SIPOptions {
	s := &SIPOptions{
		Ip:     "127.0.0.1",
		Port:   "5060",
		Domain: "4401020049",
		Id:     "44010200492000000001",
	}
	_ = viper.UnmarshalKey("sip", s)
	s.Id = "44010200492000000001"
	return s
}

func SIPAddress() string {
	return o.SIPOptions.Ip
}

func SIPPort() string {
	return o.SIPOptions.Port
}

func SIPId() string {
	return o.SIPOptions.Id
}

func SIPDomain() string {
	return o.SIPOptions.Domain
}

func SIPUserAgent() string {
	return o.SIPOptions.UserAgent
}
