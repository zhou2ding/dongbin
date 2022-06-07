package jwtauth

import "encoding/json"

type Token struct {
	AccessCode string
	Type       string
	ExpireAt   int64
}

func (t *Token) GetAccessCode() string {
	return t.AccessCode
}

func (t *Token) GetType() string {
	return t.Type
}

func (t *Token) GetExpireAt() int64 {
	return t.ExpireAt
}

func (t *Token) EncodeToJson() ([]byte, error) {
	return json.Marshal(t)
}
