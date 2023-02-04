package config

import (
	"github.com/spf13/viper"
	"time"
)

type RedisOptions struct {
	Host               string        `json:"host,omitempty" mapstructure:"host"`
	Port               int           `json:"port,omitempty" mapstructure:"port"`
	Database           int           `json:"database,omitempty" mapstructure:"database"`
	UserName           string        `json:"username" mapstructure:"username"`
	Password           string        `json:"password,omitempty" mapstructure:"password"`
	MaxRetries         int           `json:"max-retries,omitempty" mapstructure:"max-retries"`
	PoolSize           int           `json:"pool-size,omitempty" mapstructure:"pool-size"`
	MinIdleConnections int           `json:"min-idle-connections,omitempty" mapstructure:"min-idle-connections"`
	MaxIdleConnections int           `json:"max-idle-connections,omitempty" mapstructure:"max-idle-connections"`
	ConnMaxLifetime    time.Duration `json:"conn-max-life-time,omitempty" mapstructure:"conn-max-life-time"`
}

func newRedisOptions() *RedisOptions {
	r := &RedisOptions{
		Host:               "127.0.0.1",
		Port:               6379,
		Database:           0,
		Password:           "",
		MaxRetries:         3,
		PoolSize:           50,
		MinIdleConnections: 50,
		MaxIdleConnections: 100,
		ConnMaxLifetime:    0,
	}
	err := viper.UnmarshalKey("redis", r)
	if err != nil {
		panic(err)
	}
	return r
}

func RedisOption() *RedisOptions {
	return o.RedisOptions
}
