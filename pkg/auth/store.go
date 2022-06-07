package auth

type Store interface {
	// Get get expire time by token
	Get(token string) (int64, error)
	// Set add token and set its expire time
	Set(token string, expired int64, timeout int) error
	// Del delete token
	Del(token string) error
	// Check check if the token exists
	Check(token string) (bool, error)
	// SetExpired update expire time of the token
	SetExpired(token string, timeout int) error
	// Close close store
	Close() error
}
