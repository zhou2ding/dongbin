package rpcservice

const (
	ErrUnknownReqType = 20000 + iota
	ErrServiceNotSupport
	ErrMethodNotImplement
	ErrInvalidLoginRequest
	ErrInvalidUserName
	ErrInvalidJson
	ErrInvalidStationName
)

var statusText = map[int]string{
	ErrUnknownReqType:      "Unknown Request Type",
	ErrServiceNotSupport:   "Service Not Supported",
	ErrMethodNotImplement:  "Method Not Implemented",
	ErrInvalidLoginRequest: "Invalid Login Request",
	ErrInvalidUserName:     "User Not Exist Or Duplicate",
	ErrInvalidJson:         "Json Unmarshal Error",
	ErrInvalidStationName:  "Station Name Wrong",
}

func StatusText(code int) string {
	return statusText[code]
}
