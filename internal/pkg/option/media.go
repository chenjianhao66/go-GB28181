package option

import (
	"github.com/spf13/pflag"
)

type MediaOptions struct {
	Id       string `json:"id,omitempty" mapstructure:"id"`
	Ip       string `json:"ip,omitempty" mapstructure:"ip"`
	HttpPort string `json:"http-port,omitempty" mapstructure:"http-port"`
	Secret   string `json:"secret,omitempty" mapstructure:"secret"`
}

func NewMediaOption() *MediaOptions {
	return &MediaOptions{
		Id:       "FQ3TF8yT83wh5Wvz",
		Ip:       "127.0.0.1",
		HttpPort: "8000",
		Secret:   "035c73f7-bb6b-4889-a715-d9eb2d1925cc",
	}
}

func (m *MediaOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&m.Id, "media.id", m.Id, "ZLMediaKit服务的唯一id")
	fss.StringVar(&m.Ip, "media.ip", m.Ip, "ZLMediaKit服务的ip地址")
	fss.StringVar(&m.HttpPort, "media.http-port", m.HttpPort, "ZLMediaKit服务的http端口")
	fss.StringVar(&m.Secret, "media.secret", m.Secret, "ZLMediaKit服务的api密钥")
}
