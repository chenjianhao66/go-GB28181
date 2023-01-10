package config

import "time"

// MySQLOptions 定义MySQL数据库的配置选项
type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"`
	Port                  string        `json:"port"`
	Username              string        `json:"username,omitempty"`
	Password              string        `json:"password,omitempty"`
	Database              string        `json:"database,omitempty"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty"`
	LogLevel              int           `json:"log-level,omitempty"`
}

func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1",
		Port:                  "3306",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1,
	}
}
