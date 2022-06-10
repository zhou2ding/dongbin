package store

import "blog/pkg/cache"

type Store struct {
}

type AuthToken struct {
	Expired int64
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Set(token string, expire int64, timeout int) error {
	return cache.GetCache().Set("auth_"+token, AuthToken{Expired: expire}, timeout)
}

func (s *Store) Get(token string) (int64, error) {
	auth := AuthToken{}
	return auth.Expired, cache.GetCache().Get("auth_"+token, &auth, true)
}

func (s *Store) Del(token string) error {
	return cache.GetCache().Del("auth_" + token)
}

func (s *Store) Expire(token string, timeout int) error {
	return cache.GetCache().Expire("auth_"+token, timeout)
}

func (s *Store) Check(token string) (bool, error) {
	return cache.GetCache().Exists("auth_" + token)
}

func (s *Store) Close() error {
	return nil
}
