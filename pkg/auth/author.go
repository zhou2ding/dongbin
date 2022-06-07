package auth

var gAuthor Author

type Author interface {
	GenerateToken(userID string)
}

func InitAuthor() error {
	return nil
}

func GetAuthor() Author {
	return gAuthor
}
