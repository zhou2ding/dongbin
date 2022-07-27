package rabbitmq

import "blog/pkg/mqbox"

type RabbitProducer struct {
}

func newRabbitProducer(name string, eb *mqbox.ExchangeBinds) *RabbitProducer {
	return &RabbitProducer{

	}
}
