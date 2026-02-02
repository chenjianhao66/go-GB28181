package option

import "github.com/spf13/pflag"

type NutsDBOptions struct {
	Path string `json:"path,omitempty" mapstructure:"path"`
}

func NewNutsDBOptions() *NutsDBOptions {
	return &NutsDBOptions{
		Path: "./db/nuts",
	}
}

func (n *NutsDBOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&n.Path, "nutsdb.path", n.Path, "nutsdb的存储路径")
}
