package mqbox

import (
	"blog/pkg/l"
	"blog/pkg/v"
	"github.com/pkg/errors"
)

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

func InitMqBox() error {
	mqType := v.GetViper().GetString("msg_queue.mq_type")
	if mqType == "rabbitmq" || mqType == "" {

	} else if mqType == "lightmq" {

	} else {
		return errors.New("unsupported message queue type")
	}
	return nil
}

func UnInitMqBox() {
	l.Logger().Info("UnInit MqBox")
	GetMqBoxInstance().Close()
}
