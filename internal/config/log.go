package config

import (
	"github.com/spf13/viper"
	"strings"
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

func newLogOptions() *LogOptions {
	l := &LogOptions{
		Level: "info",
	}
	_ = viper.UnmarshalKey("log", l)
	return l
}

func LogLevel() string {
	level := o.LogOptions.Level
	return strings.ToUpper(level)
}

func LogFilePath() string {
	return o.LogOptions.Path
}

func LogFileName() string {
	return o.LogOptions.File
}

func LogMaxSize() int {
	return o.LogOptions.MaxSize
}

func LogMaxBackups() int {
	return o.LogOptions.MaxBackups
}

func LogMaxAge() int {
	return o.LogOptions.MaxAge
}
