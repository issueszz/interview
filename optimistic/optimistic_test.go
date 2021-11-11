package optimistic

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type Optimistic struct {
	Id int64	`gorm:"column:id; primary_key; AUTO_INCREMENT" json:"id"`
	UserId string `gorm:"column:user_id; default:0; not null" json:"user_id"`
	Amount float32 `gorm:"column:amount; not null" json:"amount"`
	Version int64 `gorm:"column:version; default:0; not null" json:"version"`
}

func TestUpdate(t *testing.T)  {

	// 连接数据库
	dsn := "root:root123@tcp(127.0.0.1:3306)/mysql_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 迁移数据库
	db.AutoMigrate(Optimistic{})

	var optimistic Optimistic
	db.First(&optimistic, &Optimistic{Id: 1})

	// 乐观锁实现
	result := db.Model(&optimistic).Where("id = ?", optimistic.Id).Where("version = ?", optimistic.Version).
		UpdateColumns(Optimistic{Amount: optimistic.Amount+10, Version: optimistic.Version + 1})
	fmt.Printf("update %v lines\n", result.RowsAffected)
}
