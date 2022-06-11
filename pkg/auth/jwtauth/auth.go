package jwtauth

import (
	"blog/pkg/auth"
	"blog/pkg/cfg"
	"blog/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"strconv"
	"time"
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
	store auth.Store
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

func New(store auth.Store, opts ...Option) *JWTAuth {
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

func (j *JWTAuth) GenerateToken(userID string) (auth.Token, error) {
	now := time.Now()
	expire := now.Add(time.Duration(j.opts.expired) * time.Second).UnixNano()

	token := jwt.NewWithClaims(j.opts.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.UnixNano(),
		ExpiresAt: expire,
		NotBefore: now.UnixNano(),
		Subject:   userID,
	})

	accessToken, err := token.SignedString(j.opts.signingKey)
	if err != nil {
		logger.GetLogger().Error("SignedString failed", zap.Error(err))
		return nil, err
	}

	err = j.callStore(func(s auth.Store) error {
		return s.Set(userID+"_"+strconv.FormatInt(now.UnixNano(), 10), expire, j.opts.expired)
	})
	if err != nil {
		logger.GetLogger().Error("set store failed", zap.Error(err))
		return nil, err
	}

	return &jwtToken{accessToken, j.opts.tokenType, expire}, nil
}

func (j *JWTAuth) GetUserFromToken(token string) (string, error) {
	claim, err := j.parseToken(token)
	if err != nil {
		return "", err
	}
	return claim.Subject, nil
}

func (j *JWTAuth) DestroyToken(token string) error {
	claim, err := j.parseToken(token)
	if err != nil {
		return err
	}

	return j.callStore(func(s auth.Store) error {
		return s.Del(claim.Subject)
	})
}

func (j *JWTAuth) UpdateToken(token string) error {
	claim, err := j.parseToken(token)
	if err != nil {
		return err
	}

	return j.callStore(func(s auth.Store) error {
		expire := time.Now().Add(time.Duration(j.opts.expired) * time.Second).UnixNano()
		return s.Set(claim.Subject+"_"+strconv.FormatInt(claim.IssuedAt, 10), expire, j.opts.expired)
	})
}

func (j *JWTAuth) parseToken(tokenStr string) (*jwt.StandardClaims, error) {
	if tokenStr == "" {
		return nil, auth.ErrInvalidToken
	}

	token, _ := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, j.opts.keyFunc)
	if token == nil || token.Claims == nil {
		return nil, auth.ErrInvalidToken
	}
	logger.GetLogger().Info("parse token", zap.String("user name", token.Claims.(*jwt.StandardClaims).Subject))
	return token.Claims.(*jwt.StandardClaims), nil
}

func (j *JWTAuth) callStore(fn func(auth.Store) error) error {
	if s := j.store; s != nil {
		return fn(s)
	}
	return nil
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
