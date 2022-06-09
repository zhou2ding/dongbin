package cache

import (
	"blog/pkg/cfg"
	"blog/pkg/logger"
	"encoding/base64"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type Redis struct {
	cli  *redis.Client
	host string
}

func InitRedis() (Cache, error) {
	host := cfg.GetViper().GetString("cache.host")
	pwd := cfg.GetViper().GetString("cache.password")

	cli := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd,
	})
	if _, err := cli.Ping().Result(); err != nil {
		logger.GetLogger().Error("connect to redis host failed", zap.String("host", host), zap.Error(err))
		return nil, err
	}

	return &Redis{cli: cli, host: host}, nil
}

func (r *Redis) Set(key string, value interface{}, timeout int) error {
	data, err := r.encode(value)
	if err != nil {
		return err
	}
	if r.cli == nil {
		if err = r.reconnect(); err != nil {
			return err
		}
	}

	expire := time.Duration(timeout) * time.Second
	if err = r.cli.Set(key, data, expire).Err(); err != nil {
		logger.GetLogger().Error("set cache failed", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

func (r *Redis) Get(key string, to interface{}, encode bool) error {
	if r.cli == nil {
		if err := r.reconnect(); err != nil {
			return err
		}
	}

	exist, err := r.Exists(key)
	if err != nil {
		logger.GetLogger().Error("get cache exists failed", zap.String("key", key), zap.Error(err))
		return err
	}
	if !exist {
		return errors.New("cache key not exist")
	}

	cmd := r.cli.Get(key)
	if cmd == nil {
		return errors.New("get cache no result")
	}
	if cmd.Err() != nil {
		logger.GetLogger().Error("get cache failed", zap.Error(err))
		return cmd.Err()
	}

	if encode {
		if err = r.decode(cmd.Val(), to); err != nil {
			logger.GetLogger().Error("get cache decode failed", zap.String("key", key), zap.Error(err))
			return err
		}
	} else {
		if err = json.Unmarshal([]byte(cmd.Val()), to); err != nil {
			logger.GetLogger().Error("get cache unmarshal failed", zap.String("key", key), zap.Error(err))
			return err
		}
	}
	return nil
}
func (r *Redis) Del(key string) error {
	if r.cli == nil {
		if err := r.reconnect(); err != nil {
			return err
		}
	}

	if err := r.cli.Del(key).Err(); err != nil {
		logger.GetLogger().Error("del cache failed", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}
func (r *Redis) Exists(key string) (bool, error) {
	if r.cli == nil {
		if err := r.reconnect(); err != nil {
			return false, err
		}
	}

	cmd := r.cli.Exists(key)
	if err := cmd.Err(); err != nil {
		logger.GetLogger().Error("exists cache failed", zap.String("key", key), zap.Error(err))
		return false, err
	}

	return cmd.Val() > 0, nil
}
func (r *Redis) Expire(key string, timeout int) error {
	if r.cli == nil {
		if err := r.reconnect(); err != nil {
			return err
		}
	}

	if err := r.cli.Expire(key, time.Duration(timeout)*time.Second).Err(); err != nil {
		logger.GetLogger().Error("set cache expire time failed", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}
func (r *Redis) Close() {
	_ = r.cli.Close()
}

func (r *Redis) reconnect() error {
	r.cli = redis.NewClient(&redis.Options{Addr: r.host})
	_, err := r.cli.Ping().Result()
	return err
}

func (r *Redis) encode(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	encodeStr := base64.StdEncoding.EncodeToString(bytes)
	return encodeStr, nil
}

func (r *Redis) decode(data string, to interface{}) error {
	decodeStr, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(decodeStr, to)
}
