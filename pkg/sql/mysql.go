package sql

import (
	"blog/pkg/v"
	"blog/pkg/internal/gormlog"
	"blog/pkg/l"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initMySQL() (*gorm.DB, error) {
	dbName := v.GetViper().GetString("database.name")
	dbUser := v.GetViper().GetString("database.name")
	dbPwd := v.GetViper().GetString("database.name")
	dbHost := v.GetViper().GetString("database.name")
	dbPort := v.GetViper().GetString("database.name")
	dbCharset := v.GetViper().GetString("database.name")

	url := dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=" + dbCharset + "&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(url), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, Logger: gormlog.NewGormLogger(l.GetLogger())})
	if err != nil {
		return nil, err
	}

	return database, nil
}
