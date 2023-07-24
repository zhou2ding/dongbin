package mqtt

import (
	"blog/pkg/l"
	"crypto/tls"
	"crypto/x509"
	"github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

type ConConfig struct {
	Broker   string
	Topic    string
	Username string
	Password string
	Cafile   string
	Cert     string
	Key      string
}

func newMqttClient(name, username, password, clientId string, topics []string) *Client {
	recvCh := make(chan interface{}, 10000)
	var publishHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
		msg := &Message{
			Duplicate: message.Duplicate(),
			Qos:       message.Qos(),
			Retained:  message.Retained(),
			Topic:     message.Topic(),
			MessageID: message.MessageID(),
			Data:      message.Payload(),
		}
		recvCh <- msg
	}
	var connHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		for _, topic := range topics {
			token := client.Subscribe(topic, 1, publishHandler)
			if token.Error() != nil || !token.WaitTimeout(15*time.Second) {
				l.Logger().Error("subscribe topic error", zap.String("topic", topic), zap.Error(token.Error()))
			}
		}
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(name).
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
		l.Logger().Error("connect to mqtt server error", zap.Error(token.Error()))
		return nil
	}
	l.Logger().Info("connect to mqtt server!", zap.String("client id", clientId))
	return cli
}

func loadTLSConfig(config *ConConfig) *tls.Config {
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = false
	if config.Cafile != "" {
		certpool := x509.NewCertPool()
		ca, err := os.ReadFile(config.Cafile)
		if err != nil {
			log.Fatalln(err.Error())
		}
		certpool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = certpool
	}
	return &tlsConfig
}
