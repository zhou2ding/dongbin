package sql

import (
	"blog/pkg/cfg"
	"blog/pkg/internal/gormlog"
	"blog/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initMySQL() (*gorm.DB, error) {
	dbName := cfg.GetViper().GetString("database.name")
	dbUser := cfg.GetViper().GetString("database.name")
	dbPwd := cfg.GetViper().GetString("database.name")
	dbHost := cfg.GetViper().GetString("database.name")
	dbPort := cfg.GetViper().GetString("database.name")
	dbCharset := cfg.GetViper().GetString("database.name")

	url := dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=" + dbCharset + "&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(url), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, Logger: gormlog.NewGormLogger(logger.GetLogger())})
	if err != nil {
		return nil, err
	}

	return database, nil
}
