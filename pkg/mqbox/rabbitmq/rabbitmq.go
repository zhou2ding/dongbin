package rabbitmq

import (
	"blog/pkg/cfg"
	"blog/pkg/mqbox"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
)

type RabbitMq struct {
	mutex sync.RWMutex

	host string

	conn *amqp.Connection

	producers map[string]mqbox.Producer
	consumers map[string]mqbox.Consumer

	closeConnChan chan *amqp.Error //notify when connection close
}

type MqInstance struct {
}

var (
	once     sync.Once
	client   *RabbitMq
	instance *MqInstance
)

func (r *RabbitMq) Init() error {
	host := cfg.GetViper().GetString("amqp.host")
	rabbitMq := &RabbitMq{
		host:          host,
		conn:          nil,
		closeConnChan: make(chan *amqp.Error, 1),
		producers:     make(map[string]mqbox.Producer),
		consumers:     make(map[string]mqbox.Consumer),
	}
	client = rabbitMq
	return nil
}

func (r *RabbitMq) GetClient() *RabbitMq {
	return client
}

func GetMqInstance() *MqInstance {
	once.Do(func() {
		instance = &MqInstance{}
	})

	return instance
}

func (r *RabbitMq) Open() error {
	if len(r.host) == 0 {
		return fmt.Errorf("AMQP host len is 0")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	var err error
	r.conn, err = amqp.Dial(r.host)
	if err != nil {
		return fmt.Errorf("dial amqp failed")
	}
	r.conn.NotifyClose(r.closeConnChan)


	return nil
}
