package rpcservice

type RPCMsg struct {
	DstId     int
	TopicName string
	MsgValue  []byte
}

func (c *RPCMsg) GetDstId() int {
	return c.DstId
}

func (c *RPCMsg) GetTopicName() string {
	return c.TopicName
}

func (c *RPCMsg) GetMsgValue() []byte {
	return c.MsgValue
}
