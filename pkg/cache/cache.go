package cache

import (
	"blog/pkg/v"
	"blog/pkg/l"
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
	l.GetLogger().Info("init cache")
	cacheType := v.GetViper().GetString("cache.type")
	switch cacheType {
	case redisType:
	case memoryType:
	}
	return nil
}

func GetCache() Cache {
	return gCache
}
