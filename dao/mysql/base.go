package mysql

import (
	"blog/dao"
	"sync"
)

type MySQL struct {
}

var (
	mySqlInstanse *MySQL
	once          sync.Once
)

func GetDB() dao.Repository {
	once.Do(func() {
		mySqlInstanse = &MySQL{}
	})
	return mySqlInstanse
}

func (m *MySQL) Start() error {
	return nil
}

func (m *MySQL) Stop() error {
	return nil
}
