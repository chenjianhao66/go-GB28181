package option

import "github.com/spf13/pflag"

// GbOption 项目选项接口，返回每个组件的命令行标志
type GbOption interface {
	Flags() (fss *pflag.FlagSet)
}

type GenericOption interface {
	AddFlags(fss *pflag.FlagSet)
}
