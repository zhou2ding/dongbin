package sql

import (
	"blog/pkg/v"
	"gorm.io/gorm"
	"sync"
)

type DataBase struct {
	dbType int
}

var (
	db       *gorm.DB
	instance *DataBase
	once     sync.Once
)

func GetDBInstance() *DataBase {
	once.Do(func() {
		instance = &DataBase{}
	})
	return instance
}

func (d *DataBase) GetDB() *gorm.DB {
	return db
}

func InitSQL() error {
	dbType := v.GetViper().GetString("database.type")
	var err error
	if dbType == "mysql" {
		db, err = initMySQL()
	}
	return err
}
