package gbctl

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/spf13/pflag"
)

type ctlOption struct {
	Sip *option.SIPOptions `json:"sip" mapstructure:"sip"`
}

func newCTLOption() *ctlOption {
	return &ctlOption{Sip: option.NewSIPOptions()}
}

func (c *ctlOption) Flags() (fss *pflag.FlagSet) {
	fss = pflag.NewFlagSet("gbctl", pflag.ExitOnError)
	c.Sip.AddFlags(fss)
	return
}
