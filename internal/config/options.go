package config

type Options struct {
	ServerOptions *ServerOptions `json:"server"`
	MySQLOptions  *MySQLOptions  `json:"mysql"`
	SIPOptions    *SIPOptions    `json:"sip"`
}

func NewOptions() *Options {
	return &Options{
		ServerOptions: NewServerOptions(),
		MySQLOptions:  NewMySQLOptions(),
		SIPOptions:    NewSIPOptions(),
	}
}
