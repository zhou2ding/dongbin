package jwtauth

import "github.com/dgrijalva/jwt-go"

const defaultKey = "zdb564"

type JWTAuth struct {
	opts *options
}

type options struct {
	signingMethod jwt.SigningMethod
	keyFunc       jwt.Keyfunc
	signingKey    interface{}
	expired       int
	tokenType     string
}

type Option func(*options)

func SetSigningMethod(m jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = m
	}
}

func SetSigningKey(k interface{}) Option {
	return func(o *options) {
		o.signingKey = k
	}
}

func SetKeyFunc(kf jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyFunc = kf
	}
}

func SetExpired(e int) Option {
	return func(o *options) {
		o.expired = e
	}
}
