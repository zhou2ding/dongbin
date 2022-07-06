package mqbox

type Msg = Message

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

func GetMqMessageType(message *Message) string {
	messageType, ok := message.Header[MessageType]
	if !ok {
		return ""
	}

	return messageType.(string)
}

func GetMqDeviceType(message *Message) int32 {
	messageType, ok := message.Header[DeviceType]
	if !ok {
		return 0
	}

	return messageType.(int32)
}
