package auth

import "github.com/pkg/errors"

var gAuthor Author

var (
	ErrInvalidToken = errors.New("unauthorized")
	ErrExpiredToken = errors.New("expired token")
)

type Token interface {
	GetAccessCode() string
	GetType() string
	GetExpireAt() int64
	EncodeToJson() ([]byte, error)
}

type Author interface {
	GenerateToken(userID string) (Token, error)
	DestroyToken(token string) error
	UpdateToken(token string) error
	CheckToken(token string) (string, string, error)
	CheckTokenWithUpdate(token string) (string, error)
	Release() error
	GetUserFromToken(token string) (string, error)
}

