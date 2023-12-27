package sqlite

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/gbserver/storage"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/model"
	"github.com/chenjianhao66/go-GB28181/internal/pkg/option"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"sync"

	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"
)

type sqliteDatastore struct {
	db *gorm.DB
}

var (
	sqliteFactory storage.Factory
	once          sync.Once
)

func GetSqliteFactory() storage.Factory {
	var (
		err   error
		dbIns *gorm.DB
		opt   option.SqliteOptions
	)

	once.Do(func() {
		if err = viper.UnmarshalKey("sqlite", &opt); err != nil {
			log.Error("load sqlite database config file fail")
			panic(err)
		}
		dbIns, err = New(&opt)
		sqliteFactory = &sqliteDatastore{dbIns}
	})

	if sqliteFactory == nil || err != nil {
		panic(fmt.Errorf("failed to get sqlite storage fatory,  error: %w", err))
	}

	return sqliteFactory
}

// New 根据MySQL选项去构建gorm对象
func New(opts *option.SqliteOptions) (*gorm.DB, error) {
	d := fmt.Sprintf("%s/%s", opts.Path, opts.File)
	db, err := gorm.Open(sqlite.Open(d), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(model.Device{}, model.MediaDetail{}, model.Channel{})
	return db, err
}

func (s *sqliteDatastore) Devices() storage.DeviceStore {
	//TODO implement me
	panic("implement me")
}

func (s *sqliteDatastore) Media() storage.MediaStorage {
	//TODO implement me
	panic("implement me")
}

func (s *sqliteDatastore) Channel() storage.ChannelStore {
	//TODO implement me
	panic("implement me")
}
