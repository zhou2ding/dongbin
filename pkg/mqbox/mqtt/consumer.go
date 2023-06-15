package mqtt

import (
	"context"
	"github.com/eclipse/paho.mqtt.golang"
	"sync"
)

type MqttClient struct {
	name      string
	topic     string
	mtx       sync.RWMutex
	ctx       context.Context
	recvQueue chan interface{}
	sendQueue chan interface{}
	cli       mqtt.Client
}
