package store

import "blog/pkg/cache"

type Storage struct {
}

type AuthToken struct {
	Expired int64
}

func NewStore() *Storage {
	return &Storage{}
}

func (s *Storage) Set(token string, expire int64, timeout int) error {
	return cache.GetCache().Set("auth_"+token, AuthToken{Expired: expire}, timeout)
}

func (s *Storage) Get(token string) (int64, error) {
	auth := AuthToken{}
	return auth.Expired, cache.GetCache().Get("auth_"+token, &auth, true)
}

func (s *Storage) Del(token string) error {
	return cache.GetCache().Del("auth_" + token)
}

func (s *Storage) Expire(token string, timeout int) error {
	return cache.GetCache().Expire("auth_"+token, timeout)
}

func (s *Storage) Check(token string) (bool, error) {
	return cache.GetCache().Exists("auth_" + token)
}

func (s *Storage) Close() error {
	return nil
}
