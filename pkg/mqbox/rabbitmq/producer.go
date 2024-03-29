package rabbitmq

import (
	"blog/pkg/l"
	"blog/pkg/mqbox"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
	"time"
)

type RabbitProducer struct {
	name          string
	mtx           sync.RWMutex
	conn          *amqp.Connection
	ch            *amqp.Channel
	exchangeBinds *mqbox.ExchangeBinds
	closeCh       chan *amqp.Error
	status        uint8
}

func newRabbitProducer(name string, eb *mqbox.ExchangeBinds) *RabbitProducer {
	return &RabbitProducer{
		name:          name,
		exchangeBinds: eb,
		closeCh:       make(chan *amqp.Error, 1),
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
		l.Logger().Error("open channel error")
		return err
	}

	if err = r.applyExchangeBinds(r.ch, r.exchangeBinds); err != nil {
		_ = r.ch.Close()
		return err
	}
	r.ch.NotifyClose(r.closeCh)

	go r.keepalive()

	r.status = mqbox.StateOpened
	l.Logger().Info("rabbitmq open success", zap.String("name", r.name))
	return nil
}

func (r *RabbitProducer) Reopen(mq interface{}) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.status == mqbox.StateOpened {
		return nil
	}

	var ok bool
	r.conn, ok = mq.(*amqp.Connection)
	if !ok {
		return errors.New("reopen mq params error")
	}

	channel, err := r.conn.Channel()
	if err != nil {
		l.Logger().Error("reopen channel failed")
		return err
	}
	r.closeCh = make(chan *amqp.Error, 1)
	channel.NotifyClose(r.closeCh)
	r.ch = channel

	if err = r.applyExchangeBinds(r.ch, r.exchangeBinds); err != nil {
		_ = r.ch.Close()
		return err
	}

	r.status = mqbox.StateOpened
	l.Logger().Info("rabbit producer reopen success", zap.String("name", r.name))
	return nil
}

func (r *RabbitProducer) Close() {
	r.mtx.Lock()
	r.ch.Close()
	close(r.closeCh)
	r.mtx.Unlock()
}

func (r *RabbitProducer) Publish(msg *mqbox.Msg) error {
	if r.conn == nil || r.ch == nil {
		return errors.New("conn or channel is nil")
	}
	if r.exchangeBinds == nil || r.exchangeBinds.Exchange == nil || r.exchangeBinds.Bindings == nil {
		return errors.New("exchangeBinds is nil")
	}

	data := amqp.Publishing{
		Headers:         msg.Header,
		ContentEncoding: "application/json",
		DeliveryMode:    mqbox.Persistent,
		Priority:        uint8(5),
		Timestamp:       time.Now(),
		Body:            msg.Body,
	}
	return r.ch.Publish(r.exchangeBinds.Exchange.Name, r.exchangeBinds.Bindings.RouteKey, false, false, data)
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
	select {
	case err := <-r.closeCh:
		if err != nil {
			l.Logger().Error("producer channel is closed with error", zap.String("name", r.name), zap.Error(err))
		} else {
			l.Logger().Info("producer channel is closed with error", zap.String("name", r.name))
		}

		r.mtx.Lock()
		r.status = mqbox.StateReopening
		r.mtx.Unlock()

		maxRetry := 99999999
		for i := 0; i < maxRetry; i++ {
			time.Sleep(8 * time.Second)
			if r.conn == nil {
				l.Logger().Error("producer connection is nil")
				return
			}
			if r.status == mqbox.StateOpened {
				l.Logger().Info("producer is opened")
				break
			}
			if err := r.Reopen(r.conn); err != nil {
				l.Logger().Info("producer reopen failed", zap.String("name", r.name), zap.Int("times", i+1), zap.Error(err))
				continue
			}
			l.Logger().Info("producer(%s) reopen done", zap.String("name", r.name), zap.Int("times", i+1))
		}
	}
}
