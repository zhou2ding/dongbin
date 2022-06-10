package jwtauth

import "encoding/json"

type jwtToken struct {
	AccessCode string
	Type       string
	ExpireAt   int64
}

func (t *jwtToken) GetAccessCode() string {
	return t.AccessCode
}

func (t *jwtToken) GetType() string {
	return t.Type
}

func (t *jwtToken) GetExpireAt() int64 {
	return t.ExpireAt
}

func (t *jwtToken) EncodeToJson() ([]byte, error) {
	return json.Marshal(t)
}
