package mqtt

import (
	"blog/pkg/l"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Client struct {
	recvCh chan interface{}
	sendCh chan interface{}
	cli    mqtt.Client
}

type Message struct {
	Duplicate bool
	Qos       byte
	Retained  bool
	Topic     string
	MessageID uint16
	Data      []byte
}

func Start(cli *Client) {
	for {
		m := <-cli.recvCh
		msg := m.(*Message)
		l.Logger().Info("receive message from broker", zap.Any("message", msg))
		// todo 增加协程池
	}
}
