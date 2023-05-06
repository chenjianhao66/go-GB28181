package gbctl

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/spf13/pflag"
)

type ctlOption struct {
	Sip          *option.SIPOptions `json:"sip" mapstructure:"sip"`
	LogOption    *option.LogOptions `json:"log" mapstructure:"log"`
	ClientOption *ClientOptions     `json:"client" mapstructure:"client"`
}

func newCTLOption() *ctlOption {
	return &ctlOption{
		Sip:          option.NewSIPOptions(),
		LogOption:    option.NewLogOptions(),
		ClientOption: newClientOption(),
	}
}

func (c *ctlOption) Flags() (fss *pflag.FlagSet) {
	fss = pflag.NewFlagSet("gbctl", pflag.ExitOnError)
	c.Sip.AddFlags(fss)
	c.LogOption.AddFlags(fss)
	c.ClientOption.AddFlags(fss)
	return
}

// ClientOptions 客户端配置
type ClientOptions struct {
	Ip                  string  `json:"ip,omitempty" mapstructure:"ip"`
	Port                string  `json:"port,omitempty" mapstructure:"port"`
	Id                  string  `json:"id,omitempty" mapstructure:"id"`
	User                string  `json:"user,omitempty" mapstructure:"user"`
	Password            string  `json:"password,omitempty" mapstructure:"password"`
	RegisterExpire      int64   `json:"registerExpire,omitempty" mapstructure:"register-expire"`
	RegisterInterval    int64   `json:"registerInterval,omitempty" mapstructure:"register-interval"`
	HeartbeatInterval   int64   `json:"heartbeatInterval,omitempty" mapstructure:"heartbeat-interval"`
	MaxHeartbeatTimeout int     `json:"maxHeartbeatTimeout,omitempty" mapstructure:"max-heartbeat-timeout"`
	Transport           string  `json:"transport,omitempty" mapstructure:"transport"`
	Alert               Alert   `json:"alert" mapstructure:"alert"`
	Channel             Channel `json:"channel" mapstructure:"channel"`
}

func newClientOption() *ClientOptions {
	return &ClientOptions{
		Port:                "5061",
		Id:                  "44010200491118000001",
		User:                "44010200491118000001",
		Password:            "root",
		RegisterExpire:      3600,
		RegisterInterval:    60,
		HeartbeatInterval:   5,
		MaxHeartbeatTimeout: 3,
		Transport:           "tcp",
		Alert:               Alert{},
		Channel:             Channel{},
	}
}

func (c *ClientOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&c.Ip, "client.ip", c.Port, "设备ip")
	fss.StringVar(&c.Port, "client.port", c.Port, "设备port")
	fss.StringVar(&c.Id, "client.id", c.Id, "设备国标id")
	fss.StringVar(&c.User, "client.user", c.User, "设备用户名")
	fss.StringVar(&c.Password, "client.password", c.Password, "设备用户名")
	fss.Int64Var(&c.RegisterExpire, "client.register-expire", c.RegisterExpire, "设备注册过期时间/秒")
	fss.Int64Var(&c.RegisterInterval, "client.register-interval", c.RegisterInterval, "设备注册间隔，最少60秒")
	fss.Int64Var(&c.HeartbeatInterval, "client.heartbeat-interval", c.HeartbeatInterval, "设备心跳间隔/妙")
	fss.IntVar(&c.MaxHeartbeatTimeout, "client.max-heartbeat-time", c.MaxHeartbeatTimeout, "设备最大心跳超时次数")
	fss.StringVar(&c.Transport, "client.transport", c.Transport, "设备传输协议，只允许传入tcp/udp")
}

type Alert struct {
	Num  int      `json:"num,omitempty" mapstructure:"num"`
	List []string `json:"list,omitempty" mapstructure:"list"`
}

type Channel struct {
	Num  int      `json:"num,omitempty" mapstructure:"num"`
	List []string `json:"list,omitempty" mapstructure:"list"`
}
