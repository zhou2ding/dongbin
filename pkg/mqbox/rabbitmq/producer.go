package rabbitmq

import (
	"blog/pkg/l"
	"blog/pkg/mqbox"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
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

func (r *RabbitProducer) Name() string {
	return r.name
}

func (r *RabbitProducer) Open(mq interface{}) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var ok bool
	r.conn, ok = mq.(*amqp.Connection)
	if !ok {
		return errors.New("open mq params error")
	}

	var err error
	r.ch, err = r.conn.Channel()
	if err != nil {
		l.GetLogger().Error("open channel error")
		return err
	}

	if err = r.applyExchangeBinds(r.ch, r.exchangeBinds); err != nil {
		_ = r.ch.Close()
		return err
	}
	r.ch.NotifyClose(r.close)

	go r.keepalive()

	r.status = mqbox.StateOpened
	l.GetLogger().Info("rabbitmq open success", zap.String("name", r.name))
	return nil
}

func (r *RabbitProducer) applyExchangeBinds(ch *amqp.Channel, binds *mqbox.ExchangeBinds) error {
	if ch == nil || binds == nil {
		return errors.New("channel or binds is nil")
	}
	if binds.Bindings == nil || binds.Exchange == nil {
		return errors.New("bindings or exchange is nil")
	}

	ex := binds.Exchange
	if err := ch.ExchangeDeclare(ex.Name, ex.Type, ex.Durable, ex.AutoDelete, ex.Internal, ex.NoWait, ex.Args); err != nil {
		return err
	}
	return nil
}

func (r *RabbitProducer) keepalive() {

}
