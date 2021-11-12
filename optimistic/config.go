package optimistic

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Db       *gorm.DB
	MaxRetry int64
}

func NewConfig(maxRetry int64) *Config {
	dsn := "root:root123@tcp(127.0.0.1:3306)/mysql_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 迁移数据库
	db.AutoMigrate(Optimistic{})

	return &Config{
		Db:       db,
		MaxRetry: maxRetry,
	}
}
