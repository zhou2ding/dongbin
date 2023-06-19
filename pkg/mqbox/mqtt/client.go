package mqtt

import (
	"blog/pkg/l"
	"github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"time"
)

type Client struct {
	recvCh chan interface{}
	sendCh chan interface{}
	cli    mqtt.Client
}

func newMqttClient(username, password, clientId string, topics []string) *Client {
	recvCh := make(chan interface{}, 10000)
	var publishHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
		recvCh <- message
	}
	var connHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		for _, topic := range topics {
			token := client.Subscribe(topic, 1, publishHandler)
			if token.Error() != nil || !token.WaitTimeout(15*time.Second) {
				l.GetLogger().Error("subscribe topic error", zap.String("topic", topic), zap.Error(token.Error()))
			}
		}
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("").
		SetUsername(username).
		SetPassword(password).
		SetClientID(clientId).
		SetDefaultPublishHandler(publishHandler).
		SetOnConnectHandler(connHandler).
		SetConnectTimeout(60 * time.Second).
		SetConnectRetry(true).
		SetAutoReconnect(false).
		SetMaxReconnectInterval(100).
		SetKeepAlive(30)

	cli := &Client{
		recvCh: recvCh,
		sendCh: make(chan interface{}, 10000),
		cli:    mqtt.NewClient(opts),
	}
	token := cli.cli.Connect()
	if token.Error() != nil || !token.WaitTimeout(15*time.Second) {
		l.GetLogger().Error("connect to mqtt server error", zap.Error(token.Error()))
		return nil
	}
	l.GetLogger().Info("connect to mqtt server!", zap.String("client id", clientId))
	return cli
}
