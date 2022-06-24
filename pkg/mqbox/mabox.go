package mqbox

type MqBox interface {
	Init()
	Open() error
	Close()
	StartMqProducer(topic string) error
	StartMqConsumer(topic string, name string, channel chan *Message) error
	Publish(topic string, msg *Message) error
}
