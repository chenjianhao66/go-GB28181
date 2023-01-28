package mysql

import (
	"fmt"
	"github.com/chenjianhao66/go-GB28181/internal/config"
	"github.com/chenjianhao66/go-GB28181/internal/log"
	"github.com/chenjianhao66/go-GB28181/internal/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	log2 "log"
	"os"
	"sync"
	"time"
)

type datastore struct {
	db *gorm.DB
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMySQLFactory get mysql database factory
func GetMySQLFactory() (store.Factory, error) {
	log.Debug("init mysql.....")
	var (
		err          error
		dbIns        *gorm.DB
		mySQLOptions config.MySQLOptions
	)
	once.Do(func() {
		if err = viper.UnmarshalKey("mysql", &mySQLOptions); err != nil {
			logrus.Error(err)
			panic("load mysql config file fail")
		}
		dbIns, err = New(&mySQLOptions)
		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}

// New 根据MySQL选项去构建gorm对象
func New(opts *config.MySQLOptions) (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Database,
		true,
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), getConfig())

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置最多连接数
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// 设置最多可重用连接
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// 设置最多空闲连接池里的最多连接数
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}

func (d *datastore) Devices() store.DeviceStore {
	return newDevices(d)
}

// 自定义gorm配置
func getConfig() *gorm.Config {
	c := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	_default := logger.New(log2.New(os.Stdout, "\r\n", log2.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond, // 打印慢SQL
		LogLevel:      logger.Info,            // 打印级别为info
		Colorful:      true,                   // 是否为彩色输出到控制台
	})
	c.Logger = _default.LogMode(logger.Error)
	return c
}
