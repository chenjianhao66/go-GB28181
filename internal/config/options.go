package config

type Options struct {
	ServerOptions *ServerOptions `json:"server" mapstructure:"server"`
	MySQLOptions  *MySQLOptions  `json:"mysql" mapstructure:"mysql"`
	SIPOptions    *SIPOptions    `json:"sip" mapstructure:"sip"`
}

func NewOptions() *Options {
	return &Options{
		ServerOptions: NewServerOptions(),
		MySQLOptions:  NewMySQLOptions(),
		SIPOptions:    NewSIPOptions(),
	}
}
