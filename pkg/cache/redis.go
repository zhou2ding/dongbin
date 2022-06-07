package cache

import (
	"blog/pkg/cfg"
	"blog/pkg/logger"
	"encoding/base64"
	"encoding/json"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
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
		logger.GetLogger().Error("connect to redis host"+host+"failed", zap.Error(err))
		return nil, err
	}

	return &Redis{cli: cli, host: host}, nil
}

func (r *Redis) Set(key string, value interface{}, timeout int) error {
	return nil
}

func (r *Redis) Get(key string, to interface{}, encode bool) error {
	return nil
}
func (r *Redis) Del(key string) error {
	return nil
}
func (r *Redis) Exists(key string) (bool, error) {
	return false, nil
}
func (r *Redis) Expire(key string, timeout int) error {
	return nil
}
func (r *Redis) Close() {

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
