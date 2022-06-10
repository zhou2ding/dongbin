package jwtauth

import (
	"blog/pkg/auth"
	"blog/pkg/auth/jwtauth/store"
	"blog/pkg/cfg"
	"github.com/dgrijalva/jwt-go"
)

const defaultKey = "zdb564"

var (
	gAuthor        auth.Author
	defaultOptions = options{
		tokenType:     "Bearer",
		expired:       60 * 60,
		signingMethod: jwt.SigningMethodHS512,
		signingKey:    []byte(defaultKey),
		keyFunc: func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrInvalidToken
			}
			return []byte(defaultKey), nil
		},
	}
)

type JWTAuth struct {
	opts  *options
	store store.Store
}

type options struct {
	signingMethod jwt.SigningMethod
	keyFunc       jwt.Keyfunc
	signingKey    interface{}
	expired       int
	tokenType     string
}

type Option func(*options)

func InitJWTAuth() (auth.Author, error) {
	var opts []Option
	if cfg.GetViper().GetInt("auth.expire") != 0 {
		opts = append(opts, setExpired(cfg.GetViper().GetInt("auth.expire")))
	}

	if cfg.GetViper().GetString("auth.signing_key") != "" {
		opts = append(opts, setSigningKey(cfg.GetViper().GetString("auth.signing_key")))
		opts = append(opts, setKeyFunc(func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrInvalidToken
			}
			return []byte(cfg.GetViper().GetString("auth.signing_key")), nil
		}))
	}

	if cfg.GetViper().GetString("auth.signing_method") != "" {
		switch cfg.GetViper().GetString("auth.signing_method") {
		case "HS256":
			opts = append(opts, setSigningMethod(jwt.SigningMethodHS256))
		case "HS384":
			opts = append(opts, setSigningMethod(jwt.SigningMethodHS384))
		case "HS512":
			opts = append(opts, setSigningMethod(jwt.SigningMethodHS512))
		}
	}

	return gAuthor, nil
}

func New(store store.Store, opts ...Option) *JWTAuth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &JWTAuth{
		opts:  &o,
		store: store,
	}
}

func GetAuthor() auth.Author {
	return gAuthor
}

func setSigningMethod(m jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = m
	}
}

func setSigningKey(k interface{}) Option {
	return func(o *options) {
		o.signingKey = k
	}
}

func setKeyFunc(kf jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyFunc = kf
	}
}

func setExpired(e int) Option {
	return func(o *options) {
		o.expired = e
	}
}
