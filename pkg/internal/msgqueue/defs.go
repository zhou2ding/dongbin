package msgqueue

type Message struct {
	Header map[string]interface{}
	Body   []byte
}
