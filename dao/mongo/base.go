package mongo

import (
	"blog/dao"
	"sync"
)

type Mongo struct {
}

var (
	mongoInstance *Mongo
	once          sync.Once
)

func GetDB() dao.Repository {
	once.Do(func() {
		mongoInstance = &Mongo{}
	})
	return mongoInstance
}

func (m *Mongo) Start() error {
	return nil
}

func (m *Mongo) Stop() error {
	return nil
}
