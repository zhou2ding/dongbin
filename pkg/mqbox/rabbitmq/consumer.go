package rabbitmq

import (
	"blog/pkg/l"
	"blog/pkg/mqbox"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
	"time"
)

type RabbitConsumer struct {
	name          string
	topic         string
	mtx           sync.RWMutex
	conn          *amqp.Connection
	ch            *amqp.Channel
	exchangeBinds *mqbox.ExchangeBinds
	prefetch      int
	closeCh       chan *amqp.Error
	status        uint8
}

func newRabbitConsumer(topic, name string, eb *mqbox.ExchangeBinds) *RabbitConsumer {
	return &RabbitConsumer{
		name:          name,
		topic:         topic,
		exchangeBinds: eb,
		closeCh:       make(chan *amqp.Error, 1),
		status:        mqbox.StateClosed,
	}
}

func (r *RabbitConsumer) Name() string {
	return r.name
}

func (r *RabbitConsumer) Topic() string {
	return r.topic
}

func (r *RabbitConsumer) Close() {
	r.mtx.Lock()
	r.ch.Close()
	close(r.closeCh)
	r.mtx.Unlock()
}

func (r *RabbitConsumer) Reopen(mq interface{}) error {
	return nil
}

func (r *RabbitConsumer) keepalive() {
	select {
	case err := <-r.closeCh:
		if err != nil {
			l.GetLogger().Error("consumer channel is closed with error", zap.String("name", r.name), zap.Error(err))
		} else {
			l.GetLogger().Info("consumer channel is closed with error", zap.String("name", r.name))
		}

		r.mtx.Lock()
		r.status = mqbox.StateReopening
		r.mtx.Unlock()

		maxRetry := 99999999
		for i := 0; i < maxRetry; i++ {
			time.Sleep(8 * time.Second)
			if r.conn == nil {
				l.GetLogger().Error("consumer connection is nil")
				return
			}
			if r.status == mqbox.StateOpened {
				l.GetLogger().Info("consumer is opened")
				break
			}
			if err := r.Reopen(r.conn); err != nil {
				l.GetLogger().Info("consumer reopen failed", zap.String("name", r.name), zap.Int("times", i+1), zap.Error(err))
				continue
			}
			l.GetLogger().Info("consumer(%s) reopen done", zap.String("name", r.name), zap.Int("times", i+1))
		}
	}
}
