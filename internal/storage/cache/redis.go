package cache

import (
	"context"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/redis/go-redis/v9"
	"net"
)

type redisClient struct {
	rdb *redis.Client
}

func newRedis() *redisClient {
	option := config.RedisOption()
	rdb := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", option.Host, option.Port),
		Username:        option.UserName,
		Password:        option.Password,
		DB:              option.Database,
		MaxRetries:      option.MaxRetries,
		PoolSize:        option.PoolSize,
		MinIdleConns:    option.MinIdleConnections,
		MaxIdleConns:    option.MaxIdleConnections,
		ConnMaxLifetime: option.ConnMaxLifetime,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Errorf("connection to redis fail,addr: %s, err: %w",
			fmt.Sprintf("%s:%d", option.Host, option.Port),
			err,
		))
	}
	log.Infof("connection to redis success,%v:%v\n", option.Host, option.Port)
	//fmt.Printf("connection to redis success,%s:%d\n", option.Host, option.Port)
	rdb.AddHook(&redisHook{})
	return &redisClient{rdb: rdb}
}

type redisHook struct{}

func (r *redisHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

func (r *redisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		log.Debugf("execute redis command:%s", cmd.Args()...)
		_ = next(ctx, cmd)
		return nil
	}
}

func (r *redisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		log.Debugf("start execute redis pipe tx command, MULTI :\n")
		for _, c := range cmds {
			log.Debugf("%s\n", c.FullName())
		}
		log.Debugf("end execute redis pipe tx command, EXEC :\n")
		_ = next(ctx, cmds)
		return nil
	}
}
