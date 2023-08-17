package cache

import (
	"blog/pkg/l"
	"blog/pkg/v"
)

const (
	redisType  = "redis"
	memoryType = "memory"
)

var gCache Cache

type Cache interface {
	Set(key string, value interface{}, timeout int) error
	Get(key string, to interface{}, encode bool) error
	Del(key string) error
	Exists(key string) (bool, error)
	Expire(key string, timeout int) error
	Close()
}

func InitCache() error {
	l.Logger().Info("init cache")
	cacheType := v.GetViper().GetString("cache.type")
	var err error
	switch cacheType {
	case redisType:
		gCache, err = InitRedis()
		if err != nil {
			return err
		}
	case memoryType:
	}
	return nil
}

func GetCache() Cache {
	return gCache
}
