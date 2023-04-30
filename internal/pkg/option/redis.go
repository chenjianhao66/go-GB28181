package option

import (
	"github.com/spf13/pflag"
)

type RedisOptions struct {
	Host               string `json:"host,omitempty" mapstructure:"host"`
	Port               int    `json:"port,omitempty" mapstructure:"port"`
	Database           int    `json:"database,omitempty" mapstructure:"database"`
	UserName           string `json:"username" mapstructure:"username"`
	Password           string `json:"password,omitempty" mapstructure:"password"`
	MaxRetries         int    `json:"max-retries,omitempty" mapstructure:"max-retries"`
	PoolSize           int    `json:"pool-size,omitempty" mapstructure:"pool-size"`
	MinIdleConnections int    `json:"min-idle-connections,omitempty" mapstructure:"min-idle-connections"`
	MaxIdleConnections int    `json:"max-idle-connections,omitempty" mapstructure:"max-idle-connections"`
	ConnMaxLifetime    int64  `json:"conn-max-life-time,omitempty" mapstructure:"conn-max-life-time"`
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Host:               "127.0.0.1",
		Port:               6379,
		Database:           0,
		MaxRetries:         3,
		PoolSize:           50,
		MinIdleConnections: 50,
		MaxIdleConnections: 100,
		ConnMaxLifetime:    0,
	}
}

func (r *RedisOptions) AddFlags(fss *pflag.FlagSet) {
	fss.StringVar(&r.Host, "redis.host", r.Host, "redis的ip地址")
	fss.IntVar(&r.Port, "redis.port", r.Port, "redis的端口")
	fss.IntVar(&r.Database, "redis.database", r.Database, "redis所选择的库序号，0-15")
	fss.StringVar(&r.Password, "redis.password", r.Password, "redis的密码")
	fss.IntVar(&r.MaxRetries, "redis.max-retries", r.MaxRetries, "最大重连次数，默认值3，-1则为放弃重试")
	fss.IntVar(&r.PoolSize, "redis.pool-size", r.PoolSize, "默认是 runtime.GOMAXPROCS * 10")
	fss.IntVar(&r.MinIdleConnections, "redis.min-idle-connections", r.MinIdleConnections, "最小空闲连接数，默认值50")
	fss.IntVar(&r.MaxIdleConnections, "redis.max-idle-connections", r.MaxIdleConnections, "最大空闲连接数，默认值100")
	fss.Int64Var(&r.ConnMaxLifetime, "redis.conn-max-life-time", r.ConnMaxLifetime, "可以重复使用连接的最长时间，0代表不关闭空闲连接，默认值0")
}
