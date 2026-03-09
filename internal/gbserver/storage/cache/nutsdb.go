package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model/constant"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/nutsdb/nutsdb"
	"github.com/pkg/errors"
)

type nutsdbClient struct {
	db *nutsdb.DB
	m  *sync.Mutex
}

const (
	defaultBucket = "cache"
)

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
	client := &nutsdbClient{
		db: db,
		m:  &sync.Mutex{},
	}
	client.createBucket()
	return client
}

func (n *nutsdbClient) Get(key string) (any, error) {
	var result []byte
	if err := n.db.View(func(tx *nutsdb.Tx) error {
		data, err := tx.Get(defaultBucket, []byte(key))
		if err != nil {
			return err
		}
		result = data
		return nil
	}); err != nil {
		log.Error(err)
		return nil, err
	}

	return result, nil
}

func (n *nutsdbClient) Set(key string, val any) {
	b, err := json.Marshal(val)
	log.Infof("设置缓存, key: %s, value: %s", key, string(b))
	if err != nil {
		log.Error(err)
		return
	}

	if err = n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(defaultBucket, []byte(key), b, 0)
	}); err != nil {
		log.Error(err)
	}
}

func (n *nutsdbClient) Del(key string) error {
	if err := n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(defaultBucket, []byte(key))
	}); err != nil {
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
		data, err := tx.Get(defaultBucket, []byte(constant.CeqPrefix))
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
		return tx.Put(defaultBucket, []byte(constant.CeqPrefix), ceqBytes, 0)
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

func (n *nutsdbClient) createBucket() {
	if err := n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.NewBucket(nutsdb.DataStructureBTree, defaultBucket)
	}); err != nil {
		if errors.Is(err, nutsdb.ErrBucketAlreadyExist) {
			log.Info("default bucket已存在, 跳过创建")
			return
		}
		log.Error("创建default bucket失败")
	}

	// 初始化序号
	if err := n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(defaultBucket, []byte(constant.CeqPrefix), []byte("1"), 0)
	}); err != nil {
		log.Error("初始化序号失败", err)
	}
}
