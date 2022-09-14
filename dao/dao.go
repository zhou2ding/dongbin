package dao

type Repository interface {
	Start() error
	Stop() error
}
