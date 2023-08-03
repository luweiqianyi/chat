package mysql

import (
	"chat/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

var once sync.Once
var gClient *Client

type Config struct {
	DSN string
}

type Client struct {
	db *gorm.DB
}

func NewClient(cfg Config) *Client {
	once.Do(func() {
		db := InitDB(cfg.DSN)
		gClient = &Client{
			db: db,
		}
	})
	return gClient
}

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别，该级别下可以查看执行的SQL语句
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func (cli *Client) GetDB() *gorm.DB {
	if cli.db == nil {
		log.Panicf("db required")
		return nil
	}
	return cli.db
}
