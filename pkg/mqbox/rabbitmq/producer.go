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
		l.GetLogger().Error("reopen channel failed")
		return err
	}
	r.close = make(chan *amqp.Error, 1)
	channel.NotifyClose(r.close)
	r.ch = channel

	if err = r.applyExchangeBinds(r.ch, r.exchangeBinds); err != nil {
		_ = r.ch.Close()
		return err
	}

	r.status = mqbox.StateOpened
	l.GetLogger().Info("rabbit producer reopen success", zap.String("name", r.name))
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
	select {
	case err := <-r.close:
		if err != nil {
			l.GetLogger().Error("producer channel is closed with error", zap.String("name", r.name), zap.Error(err))
		} else {
			l.GetLogger().Info("producer channel is closed with error", zap.String("name", r.name))
		}

		r.mtx.Lock()
		r.status = mqbox.StateReopening
		r.mtx.Unlock()

		maxRetry := 99999999
		for i := 0; i < maxRetry; i++ {
			time.Sleep(8 * time.Second)
			if r.conn == nil {
				l.GetLogger().Error("producer connection is nil")
				return
			}
			if r.status == mqbox.StateOpened {
				l.GetLogger().Info("producer is opened")
				break
			}
			if err := r.Reopen(r.conn); err != nil {
				l.GetLogger().Info("producer reopen failed", zap.String("name", r.name), zap.Int("times", i+1), zap.Error(err))
				continue
			}
			l.GetLogger().Info("producer(%s) reopen done", zap.String("name", r.name), zap.Int("times", i+1))
		}
	}
}
