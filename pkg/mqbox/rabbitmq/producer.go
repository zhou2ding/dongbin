package rabbitmq

import (
	"blog/pkg/mqbox"
	"github.com/streadway/amqp"
	"sync"
)

type RabbitProducer struct {
	name          string
	mtx           sync.RWMutex
	conn          *amqp.Connection
	ch            *amqp.Channel
	exchangeBinds *mqbox.ExchangeBinds
	close         chan *amqp.Error
	status        uint8
}

func newRabbitProducer(name string, eb *mqbox.ExchangeBinds) *RabbitProducer {
	return &RabbitProducer{
		name:          name,
		exchangeBinds: eb,
		close:         make(chan *amqp.Error, 1),
		status:        mqbox.StateClosed,
	}
}
