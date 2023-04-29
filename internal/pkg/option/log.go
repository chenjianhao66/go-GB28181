package option

import (
	"github.com/spf13/pflag"
)

// TODO 添加日志配置

type LogOptions struct {
	Level      string `json:"level" mapstructure:"level"`
	Path       string `json:"path" mapstructure:"path"`
	File       string `json:"file" mapstructure:"file"`
	MaxSize    int    `json:"maxSize" mapstructure:"maxSize"`
	MaxBackups int    `json:"maxBackups" mapstructure:"maxBackups"`
	MaxAge     int    `json:"maxAge" mapstructure:"maxAge"`
}

func NewLogOptions() *LogOptions {
	return &LogOptions{
		Level:      "info",
		Path:       "./log",
		File:       "gp-gb28181.log",
		MaxSize:    1,
		MaxBackups: 30,
		MaxAge:     30,
	}
}

// AddFlags 把log的参数添加到传入的命令行标志集中
func (l *LogOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&l.Level, "log.level", l.Level, "服务的日志级别")
	fss.StringVar(&l.Path, "log.path", l.Path, "服务日志的输出目录")
	fss.StringVar(&l.File, "log.file", l.File, "服务日志输出的文件名称")
	fss.IntVar(&l.MaxSize, "log.max-size", l.MaxSize, "日志文件在轮换之前的最大大小，M为单位")
	fss.IntVar(&l.MaxBackups, "log.max-backups", l.MaxBackups, "是要保留的最大旧日志文件数。默认是保留所有旧的日志文件")
	fss.IntVar(&l.MaxAge, "log.max-age", l.MaxAge, "是根据文件名中编码的时间戳保留旧日志文件的最大天数")
}
