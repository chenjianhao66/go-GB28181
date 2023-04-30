package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/model/constant"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"net"
	"sync"
	"time"
)

type redisClient struct {
	rdb *redis.Client
	m   *sync.Mutex
}

func newRedis(opt *option.RedisOptions) *redisClient {
	//if err := viper.UnmarshalKey("redis", opt); err != nil {
	//	_, _ = fmt.Fprintf(os.Stderr, "解析redis配置失败,err : %v", err)
	//	os.Exit(1)
	//}
	rdb := redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		Username:        opt.UserName,
		Password:        opt.Password,
		DB:              opt.Database,
		MaxRetries:      opt.MaxRetries,
		PoolSize:        opt.PoolSize,
		MinIdleConns:    opt.MinIdleConnections,
		MaxIdleConns:    opt.MaxIdleConnections,
		ConnMaxLifetime: time.Duration(opt.ConnMaxLifetime),
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Errorf("connection to redis fail,addr: %s, err: %w",
			fmt.Sprintf("%s:%d", opt.Host, opt.Port),
			err,
		))
	}
	log.Infof("connection to redis success,%v:%v\n", opt.Host, opt.Port)
	//fmt.Printf("connection to redis success,%s:%d\n", options.Host, options.Port)
	rdb.AddHook(&redisHook{})
	return &redisClient{
		rdb: rdb,
		m:   &sync.Mutex{},
	}
}

func (r *redisClient) Get(key string) (any, error) {
	result, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		log.Error(err)
	}
	return result, err
}

func (r *redisClient) Set(key string, val any) {
	b, _ := json.MarshalIndent(val, "", "  ")
	if err := r.rdb.Set(context.Background(), key, b, redis.KeepTTL).Err(); err != nil {
		log.Error(err)
	}
}

func (r *redisClient) Del(key string) error {
	_, err := r.rdb.Del(context.Background(), key).Result()
	if err != nil {
		log.Error(err)
		return errors.New(err.Error())
	}
	return err
}

func (r *redisClient) GetCeq() (int64, error) {
	r.m.Lock()
	defer r.m.Unlock()
	return r.rdb.Incr(context.Background(), constant.CeqPrefix).Result()
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
