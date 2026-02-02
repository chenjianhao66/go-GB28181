package gbserver

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/spf13/pflag"
)

type GbOption struct {
	ServerOption *option.ServerOptions `json:"server,omitempty" mapstructure:"server"`
	MediaOption  *option.MediaOptions  `json:"media,omitempty" mapstructure:"media"`
	SqliteOption *option.SqliteOptions `json:"sqliteOption" mapstructure:"sqlite"`
	NutsDBOption *option.NutsDBOptions `json:"nutsdb,omitempty" mapstructure:"nutsdb"`
	LogOption    *option.LogOptions    `json:"log,omitempty" mapstructure:"log"`
	Sip          *option.SIPOptions    `json:"sip" mapstructure:"sip"`
}

func newGbOption() *GbOption {
	return &GbOption{
		ServerOption: option.NewServerOptions(),
		MediaOption:  option.NewMediaOption(),
		SqliteOption: option.NewSqliteOptions(),
		NutsDBOption: option.NewNutsDBOptions(),
		LogOption:    option.NewLogOptions(),
		Sip:          option.NewSIPOptions(),
	}
}

func (c *GbOption) Flags() (fss *pflag.FlagSet) {
	fss = pflag.NewFlagSet("gbserver", pflag.ExitOnError)
	c.ServerOption.AddFlags(fss)
	c.SqliteOption.AddFlags(fss)
	c.MediaOption.AddFlags(fss)
	c.NutsDBOption.AddFlags(fss)
	c.LogOption.AddFlags(fss)
	c.Sip.AddFlags(fss)
	return
}
