package mqtt

import (
	"blog/pkg/l"
	"blog/pkg/mqbox"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Client struct {
	recvCh chan interface{}
	sendCh chan interface{}
	cli    mqtt.Client
	name   string
	topics []string
}

type Message struct {
	Duplicate bool
	Qos       byte
	Retained  bool
	Topic     string
	MessageID uint16
	Data      []byte
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Topic() []string {
	return c.topics
}

func (c *Client) Open(mq interface{}) error {
	for {
		m := <-c.recvCh
		msg, ok := m.(*Message)
		if !ok {
			return nil
		}
		l.Logger().Info("receive message from broker", zap.Any("message", msg))
		// todo 增加协程池
	}
}

func (c *Client) Reopen(mq interface{}) error {
	return nil
}

func (c *Client) Close() {
	c.cli.Disconnect(100)
}

func (c *Client) SetMsgCallback(cb chan<- *mqbox.Message) {

}
func (c *Client) StartConsume() error {
	return nil
}

func (c *Client) Publish(msg *Message) error {
	return nil
}
