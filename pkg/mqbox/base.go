package mqbox

import "blog/pkg/internal/msgqueue"

type Message = msgqueue.Message

const (
	MessageType = "MessageType"
	DeviceType  = "DeviceType"
)

const (
	MqHeaderAlertRule      = "alertRule"
	MqHeaderDetectionData  = "detectionData"
	MqHeaderSimulationData = "simulationData"
)

const (
	MqNameEvent     = "event"
	MqNameDetection = "detection"
)

const (
	DevicePanto int32 = 1
	Device360   int32 = 2
	DeviceTwm   int32 = 4
	DeviceTwd2D int32 = 8
	DeviceTwd3D int32 = 16
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
