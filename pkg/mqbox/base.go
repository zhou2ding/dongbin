package mqbox

import "blog/pkg/internal/msgqueue"

type Message = msgqueue.Message

const (
	MessageType = "MessageType"
	DeviceType  = "DeviceType"
)

const (
	MqHeaderDemo = "demo_header"
)

const (
	MqNameDemo = "demo_name"
)

func GetMqMessageType(message *msgqueue.Message) string {
	messageType, ok := message.Header[MessageType]
	if !ok {
		return ""
	}

	return messageType.(string)
}

func GetMqDeviceType(message *msgqueue.Message) int32 {
	messageType, ok := message.Header[DeviceType]
	if !ok {
		return 0
	}

	return messageType.(int32)
}
