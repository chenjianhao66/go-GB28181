package config

type SIPOptions struct {
	Ip       string `json:"ip,omitempty"`
	Port     string `json:"port,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Id       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

func NewSIPOptions() *SIPOptions {
	return &SIPOptions{
		Ip:     "127.0.0.1",
		Port:   "5060",
		Domain: "4401020049",
		Id:     "44010200492000000001",
	}
}
