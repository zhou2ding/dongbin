package msgqueue

type Consumer interface {
	Name() string
	Topic() string

	Open(mq interface{}) error
	Reopen(mq interface{}) error
	Close()

	SetMsgCallback(cb chan<- *Message)
	StartConsume() error
}
