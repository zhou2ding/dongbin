package mqbox

import (
	"github.com/streadway/amqp"
	"time"
)

type Message struct {
	Header map[string]interface{}
	Body   []byte
}

// ExchangeBinds exchange ==> routeKey ==> queues
type ExchangeBinds struct {
	Exchange *Exchange
	Bindings *Binding
}

// Binding routeKey ==> queues
type Binding struct {
	RouteKey string
	Queues   *Queue
	NoWait   bool       // default is false
	Args     amqp.Table // default is nil
}

// Exchange 基于amqp的Exchange配置
type Exchange struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table // default is nil
}

// Exchange type
var (
	ExchangeDirect  = amqp.ExchangeDirect
	ExchangeFanout  = amqp.ExchangeFanout
	ExchangeTopic   = amqp.ExchangeTopic
	ExchangeHeaders = amqp.ExchangeHeaders
)

// DeliveryMode
var (
	Transient  uint8 = amqp.Transient
	Persistent uint8 = amqp.Persistent
)

var (
	StateOpened    = uint8(0)
	StateClosed    = uint8(1)
	StateReopening = uint8(2)
)

func DefaultExchange(name string, kind string) *Exchange {
	return &Exchange{
		Name:       name,
		Type:       kind,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}

// Queue 基于amqp的Queue配置
type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func DefaultQueue(name string) *Queue {
	return &Queue{
		Name:       name,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

// PublishMsg 生产者生产的数据格式
type PublishMsg struct {
	ContentType     string // MIME content type
	ContentEncoding string // MIME content type
	DeliveryMode    uint8  // Transient or Persistent
	Priority        uint8  // 0 to 9
	Timestamp       time.Time
	Body            []byte
}

// ConsumeOption 消费者消费选项
type ConsumeOption struct {
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

func DefaultConsumeOption() *ConsumeOption {
	return &ConsumeOption{
		NoWait:  true,
		AutoAck: true,
	}
}
