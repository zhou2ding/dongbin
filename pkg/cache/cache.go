type Cache interface {
	SetChache(key string, value interface{}, timeOut int) error
	GetChache(key string) (interface{}, error)
}