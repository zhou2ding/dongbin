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

type RabbitConsumer struct {
	name          string
	topic         string
	mtx           sync.RWMutex
	conn          *amqp.Connection
	ch            *amqp.Channel
	exchangeBinds *mqbox.ExchangeBinds
	prefetch      int
	closeCh       chan *amqp.Error
	stopChan      chan struct{}
	status        uint8
	callback      chan<- *mqbox.Message
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

func (r *RabbitConsumer) Topic() []string {
	return []string{r.topic}
}

func (r *RabbitConsumer) Open(mq interface{}) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var ok bool
	r.conn, ok = mq.(*amqp.Connection)
	if !ok {
		return errors.New("open param error")
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

func (r *RabbitConsumer) SetQos(prefetch int) {
	r.prefetch = prefetch
}

func (r *RabbitConsumer) SetMsgCallback(callback chan<- *mqbox.Msg) {
	r.mtx.Lock()
	r.callback = callback
	r.mtx.Unlock()
}

func (r *RabbitConsumer) StartConsume() error {
	if r.exchangeBinds == nil || r.exchangeBinds.Bindings == nil || r.exchangeBinds.Bindings.Queues == nil {
		return errors.New("exchangeBinds is nil")
	}

	opt := mqbox.DefaultConsumeOption()
	delivery, err := r.ch.Consume(r.exchangeBinds.Bindings.Queues.Name, "", opt.AutoAck, opt.Exclusive, opt.NoLocal, opt.NoWait, opt.Args)
	if err != nil {
		return err
	}
	if delivery == nil {
		return errors.New("consume no result")
	}
	go r.deliver(delivery)
	return nil
}

func (r *RabbitConsumer) deliver(delivery <-chan amqp.Delivery) {
	for {
		select {
		case d := <-delivery:
			if r.callback != nil && d.Body != nil {
				msg := mqbox.Msg{
					Header: d.Headers,
					Body:   d.Body,
				}
				r.callback <- &msg
			}
		case <-r.stopChan:
			l.Logger().Info("receive stop signal, stop deliver")
			return
		}
	}
}

func (r *RabbitConsumer) Close() {
	r.mtx.Lock()
	r.ch.Close()
	close(r.closeCh)
	r.mtx.Unlock()
}

func (r *RabbitConsumer) Reopen(mq interface{}) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if r.status == mqbox.StateOpened {
		return nil
	}

	var ok bool
	r.conn, ok = mq.(*amqp.Connection)
	if !ok {
		return errors.New("reopen param error")
	}

	channel, err := r.conn.Channel()
	if err != nil {
		return err
	}
	close(r.stopChan)
	time.Sleep(time.Millisecond * 200)

	r.stopChan = make(chan struct{})
	r.closeCh = make(chan *amqp.Error, 1)
	channel.NotifyClose(r.closeCh)
	r.ch = channel

	err = func(ch *amqp.Channel) error {
		if err := r.applyExchangeBinds(ch, r.exchangeBinds); err != nil {
			return err
		}
		if err := ch.Qos(r.prefetch, 0, false); err != nil {
			return err
		}
		return nil
	}(channel)

	r.status = mqbox.StateOpened
	l.Logger().Info("rabbitmq reopen success", zap.String("name", r.name))
	return nil
}

func (r *RabbitConsumer) applyExchangeBinds(ch *amqp.Channel, binds *mqbox.ExchangeBinds) error {
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

func (r *RabbitConsumer) keepalive() {
	select {
	case err := <-r.closeCh:
		if err != nil {
			l.Logger().Error("consumer channel is closed with error", zap.String("name", r.name), zap.Error(err))
		} else {
			l.Logger().Info("consumer channel is closed with error", zap.String("name", r.name))
		}

		r.mtx.Lock()
		r.status = mqbox.StateReopening
		r.mtx.Unlock()

		maxRetry := 99999999
		for i := 0; i < maxRetry; i++ {
			time.Sleep(8 * time.Second)
			if r.conn == nil {
				l.Logger().Error("consumer connection is nil")
				return
			}
			if r.status == mqbox.StateOpened {
				l.Logger().Info("consumer is opened")
				break
			}
			if err := r.Reopen(r.conn); err != nil {
				l.Logger().Info("consumer reopen failed", zap.String("name", r.name), zap.Int("times", i+1), zap.Error(err))
				continue
			}
			l.Logger().Info("consumer(%s) reopen done", zap.String("name", r.name), zap.Int("times", i+1))
		}
	}
}
