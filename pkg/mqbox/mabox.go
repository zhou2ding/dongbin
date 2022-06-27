package mqbox

type MqBox interface {
	Init() error
	Open() error
	Close()
	StartMqProducer(topic string) error
	StartMqConsumer(topic string, name string, channel chan *Message) error
	Publish(topic string, msg *Message) error
}

var gMqBox MqBox

func GetMqBoxInstance() MqBox {
	return gMqBox
}
