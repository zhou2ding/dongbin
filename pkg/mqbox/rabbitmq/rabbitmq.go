package rabbitmq

import (
	"blog/pkg/l"
	"blog/pkg/mqbox"
	"blog/pkg/v"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"sync"
	"time"
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
	host := v.GetViper().GetString("amqp.host")
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

	go r.keepalive()

	return nil
}

func (r *RabbitMq) Close() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, p := range r.producers {
		p.Close()
	}
	r.producers = make(map[string]mqbox.Producer)

	for _, p := range r.consumers {
		p.Close()
	}
	r.consumers = make(map[string]mqbox.Consumer)

	if r.conn != nil {
		_ = r.conn.Close()
	}
}

func (r *RabbitMq) StartMqProducer(topic string) error {
	ex := &mqbox.ExchangeBinds{
		Exchange: mqbox.DefaultExchange("exch."+topic, mqbox.ExchangeFanout),
		Bindings: &mqbox.Binding{
			RouteKey: "router." + topic,
			Queues:   mqbox.DefaultQueue("queue." + topic),
		},
	}

	newRabbitProducer(topic, ex)

	return nil
}

func (r *RabbitMq) StartMqConsumer(topic string, name string, channel chan *mqbox.Message) error {
	return nil
}

func (r *RabbitMq) keepalive() {
	select {
	case err := <-r.closeConnChan:
		if err != nil {
			l.Logger().Error("AMQP connection was closed with error", zap.Error(err))
		} else {
			l.Logger().Error("AMQP connection was closed with no error")
		}
		maxRetry := 99999999
		for i := 0; i < maxRetry; i++ {
			time.Sleep(5 * time.Second)
			if err2 := r.reopen(); err2 != nil {
				l.Logger().Info("AMQP reconnect failed", zap.Int("retry times", i+1), zap.Error(err2))
				continue
			}

			for _, v := range r.producers {
				e := v.Reopen(r.conn)
				if e != nil {
					l.Logger().Info("producer reopen failed", zap.Error(e))
				}
			}

			for _, v := range r.consumers {
				e := v.Reopen(r.conn)
				if e != nil {
					l.Logger().Info("consumer reopen failed", zap.Error(e))
				} else {
					v.StartConsume()
				}
			}
		}
	}
}

func (r *RabbitMq) reopen() error {
	if len(r.host) == 0 {
		l.Logger().Info("AMQP host len is 0")
		return fmt.Errorf("AMQP host len is 0")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	var err error
	r.conn, err = amqp.Dial(r.host)
	if err != nil {
		l.Logger().Info("dial amqp failed")
		return err
	}

	r.closeConnChan = make(chan *amqp.Error, 1)
	r.conn.NotifyClose(r.closeConnChan)

	return nil
}
