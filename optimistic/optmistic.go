package optimistic

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type CallBack func(Lock) Lock

var callBack = func(lock Lock) Lock {
	bizModel := lock.(*Optimistic)
	bizModel.Amount += 10
	return bizModel
}

type Lock interface {
	GetVersion() int64
	SetVersion(version int64)
}

func (o *Optimistic) GetVersion() int64 {
	return o.Version
}

func (o *Optimistic) SetVersion(version int64) {
	o.Version = version
}

// UpdateWithOptimistic 带乐观锁的mysql update
func UpdateWithOptimistic(db *gorm.DB, model Lock, callBack CallBack, currentRetry, maxRetry int64) error {
	// 递归截止
	if currentRetry > maxRetry {
		return errors.New("maximum number of executions")
	}

	currentVersion := model.GetVersion()
	model.SetVersion(currentVersion + 1)

	// 更新数据
	column := db.Model(model).Where("version = ?", currentVersion).UpdateColumns(model)

	// 更新失败
	if column.RowsAffected == 0 {
		if callBack == nil && maxRetry == 0 {
			return errors.New("concurrent optimistic update failed")
		}

		time.Sleep(time.Microsecond)
		db.First(model)
		newModel := callBack(model)
		currentRetry++

		return UpdateWithOptimistic(db, newModel, callBack, currentRetry, maxRetry)
	}
	return column.Error
}

func Update() error {
	config := NewConfig(3)

	var optimistic Optimistic
	config.Db.First(&optimistic)
	optimistic.Amount += 10

	return UpdateWithOptimistic(config.Db, &optimistic, callBack, 0, config.MaxRetry)
}
