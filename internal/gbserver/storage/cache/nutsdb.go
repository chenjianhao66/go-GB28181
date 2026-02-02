package cache

import (
	"encoding/json"
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model/constant"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/nutsdb/nutsdb"
	"os"
	"sync"
)

type nutsdbClient struct {
	db *nutsdb.DB
	m  *sync.Mutex
}

func newNutsDB(opt *option.NutsDBOptions) *nutsdbClient {
	// 确保目录存在
	if err := os.MkdirAll(opt.Path, os.ModePerm); err != nil {
		panic(fmt.Errorf("创建nutsdb目录失败: %w", err))
	}

	opts := nutsdb.DefaultOptions
	opts.Dir = opt.Path
	db, err := nutsdb.Open(opts)
	if err != nil {
		panic(fmt.Errorf("连接nutsdb失败, path: %s, err: %w", opt.Path, err))
	}

	log.Infof("连接nutsdb成功, path: %s", opt.Path)
	return &nutsdbClient{
		db: db,
		m:  &sync.Mutex{},
	}
}

func (n *nutsdbClient) Get(key string) (any, error) {
	var result []byte
	err := n.db.View(func(tx *nutsdb.Tx) error {
		data, err := tx.Get("cache", []byte(key))
		if err != nil {
			return err
		}
		result = data
		return nil
	})

	if err != nil {
		log.Error(err)
		return nil, err
	}

	// 尝试解析为JSON
	var val any
	if err := json.Unmarshal(result, &val); err != nil {
		// 如果不是JSON，直接返回字符串
		return string(result), nil
	}

	return val, nil
}

func (n *nutsdbClient) Set(key string, val any) {
	b, err := json.Marshal(val)
	if err != nil {
		log.Error(err)
		return
	}

	err = n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("cache", []byte(key), b, 0)
	})

	if err != nil {
		log.Error(err)
	}
}

func (n *nutsdbClient) Del(key string) error {
	err := n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete("cache", []byte(key))
	})

	if err != nil {
		log.Error(err)
		return fmt.Errorf("删除key失败: %w", err)
	}

	return nil
}

func (n *nutsdbClient) GetCeq() (int64, error) {
	n.m.Lock()
	defer n.m.Unlock()

	var ceq int64 = 1
	err := n.db.View(func(tx *nutsdb.Tx) error {
		data, err := tx.Get("counter", []byte(constant.CeqPrefix))
		if err != nil {
			return err
		}

		// 简单的字节数组转换
		for _, b := range data {
			ceq = ceq*10 + int64(b)
		}
		return nil
	})

	if err != nil {
		log.Error(err)
		return 0, err
	}

	// 递增并保存
	ceq++
	err = n.db.Update(func(tx *nutsdb.Tx) error {
		ceqBytes := []byte(fmt.Sprintf("%d", ceq))
		return tx.Put("counter", []byte(constant.CeqPrefix), ceqBytes, 0)
	})

	if err != nil {
		log.Error(err)
		return 0, err
	}

	return ceq, nil
}

func (n *nutsdbClient) Close() error {
	if n.db != nil {
		return n.db.Close()
	}
	return nil
}
