package gbctl

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/spf13/pflag"
)

type ctlOption struct {
	Sip       *option.SIPOptions `json:"sip" mapstructure:"sip"`
	LogOption *option.LogOptions `json:"log" mapstructure:"log"`
}

func newCTLOption() *ctlOption {
	return &ctlOption{
		Sip:       option.NewSIPOptions(),
		LogOption: option.NewLogOptions(),
	}
}

func (c *ctlOption) Flags() (fss *pflag.FlagSet) {
	fss = pflag.NewFlagSet("gbctl", pflag.ExitOnError)
	c.Sip.AddFlags(fss)
	c.LogOption.AddFlags(fss)
	return
}
