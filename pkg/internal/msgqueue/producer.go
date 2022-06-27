package msgqueue

type Producer interface {
	Name() string

	Open(mq interface{}) error
	Reopen(mq interface{}) error
	Close()

	Publish(msg *Message) error
}
