package option

import (
	"github.com/spf13/pflag"
)

// MySQLOptions 定义MySQL数据库的配置选项
type MySQLOptions struct {
	Host                  string `json:"host,omitempty" mapstructure:"host"`
	Port                  string `json:"port" mapstructure:"port"`
	Username              string `json:"username,omitempty" mapstructure:"username"`
	Password              string `json:"password,omitempty" mapstructure:"password"`
	Database              string `json:"database,omitempty" mapstructure:"database"`
	MaxIdleConnections    int    `json:"max-idle-connections,omitempty" mapstructure:"max-idle-connections"`
	MaxOpenConnections    int    `json:"max-open-connections,omitempty" mapstructure:"max-open-connections"`
	MaxConnectionLifeTime int64  `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int    `json:"log-level,omitempty" mapstructure:"log-level"`
}

func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1",
		Port:                  "3306",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: 10,
		LogLevel:              1,
	}
}

func (m *MySQLOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&m.Host, "mysql.host", m.Host, "mysql数据库的ip地址")
	fss.StringVar(&m.Port, "mysql.port", m.Port, "mysql数据库的端口")
	fss.IntVar(&m.MaxIdleConnections, "mysql.max-idle-connections", m.MaxIdleConnections, "mysql数据库的最大空闲连接数")
	fss.IntVar(&m.MaxOpenConnections, "mysql.max-open-connections", m.MaxOpenConnections, "mysql数据库的最大连接数")
	fss.Int64Var(&m.MaxConnectionLifeTime, "mysql.max-connection-life-time", m.MaxConnectionLifeTime, "mysql数据库的最大可重用连接数")
	fss.IntVar(&m.LogLevel, "mysql.log-level", m.LogLevel, "mysql数据库的sql日志打印级别")
}
